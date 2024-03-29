// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/juju/errgo"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	ansi "github.com/makyo/ansigo"
	"github.com/makyo/gotui"

	"github.com/makyo/stimmtausch/client"
	"github.com/makyo/stimmtausch/help"
	"github.com/makyo/stimmtausch/signal"
	"github.com/makyo/stimmtausch/util"
)

type tui struct {
	g             *gotui.Gui
	client        *client.Client
	sent          *History
	errs          *History
	currView      *receivedView
	currViewIndex int
	views         []*receivedView
	listener      chan signal.Signal
	ready         chan bool
	title         string
	titleLen      int
	modalOpen     bool
}

var log = loggo.GetLogger("stimmtausch.ui")

// view returns the view for the given name, similar to gotui.Gui.View()
func (t *tui) view(name string) (*receivedView, error) {
	for _, v := range t.views {
		if v.viewName == name {
			return v, nil
		}
	}
	return nil, errgo.Newf("receivedView %s not found", name)
}

// connect tells the client to connect to the provided connection string, If
// successful, it will construct a receivedView to represent and hold that
// connection.
func (t *tui) connect(name string, g *gotui.Gui) error {
	log.Tracef("creating a connection with connection string %s", name)
	conn, ok := t.client.Conn(name)
	if !ok {
		log.Errorf("unable to find connection %s", name)
		return nil
	}
	connName := conn.GetConnectionName()

	for _, v := range t.views {
		if connName == v.connName {
			v.conn = conn
			conn.AddOutput(v.viewName, v.buffer, true)
			log.Tracef("opening connection for %s", name)
			err := conn.Open()
			if err != nil {
				log.Errorf("unable to open connection for %s: %v", name, err)
				return errgo.Mask(err)
			}
			t.currView = v
			t.currViewIndex = v.index
			t.currView.connected = conn.Connected
			t.updateSendTitle()
			return nil
		}
	}

	viewName := fmt.Sprintf("recv%d", len(t.views))
	log.Tracef("building received view %s", viewName)
	if v, err := g.SetView(viewName, -3, -3, -1, -1); err != nil {
		if err != gotui.ErrUnknownView {
			log.Warningf("unable to create view %+v", err)
			return errgo.Mask(err)
		}

		// Prime the view with a newline, which keeps it from complaining about
		// coordinates later.
		fmt.Fprintln(v, util.Logo())
		fmt.Fprintln(v, ansi.MaybeApply("bold+243", "\nConnecting to "+conn.GetDisplayName()+"...\n"))
		v.Wrap = true
		v.WordWrap = true
		v.IndentFirst = t.client.Config.Client.UI.IndentFirst
		v.IndentSubsequent = t.client.Config.Client.UI.IndentSubsequent
		v.Frame = false
		t.currView = &receivedView{
			connName:    connName,
			displayName: conn.GetDisplayName(),
			conn:        conn,
			viewName:    viewName,
			buffer:      NewHistory(t.client.Config.Client.UI.Scrollback, false),
			maxBuffer:   conn.GetMaxBuffer(),
			current:     true,
			index:       len(t.views),
		}
		for _, v := range t.views {
			v.current = false
		}
		t.views = append(t.views, t.currView)
		t.updateSendTitle()

		// Set this view's index as the current view index.
		t.currViewIndex = t.currView.index

		// Attach a hook that writes to the view when a line is received in
		// the received history for the connection.
		tv := t.currView
		t.currView.buffer.AddPostWriteHook(func(line *HistoryLine) error {
			fmt.Fprint(v, line.Text)
			g.Update(func(gg *gotui.Gui) error {
				return errgo.Mask(tv.updateRecvOrigin(t.currViewIndex, gg, t))
			})
			return nil
		})

		// Add the received history to the connection as an output.
		conn.AddOutput(viewName, t.currView.buffer, true)

		log.Tracef("opening connection for %s", name)
		err = conn.Open()
		if err != nil {
			log.Errorf("unable to open connection for %s: %v", name, err)
			return errgo.Mask(err)
		}
		t.currView.connected = conn.Connected
		t.updateSendTitle()
	}
	return nil
}

// postCreate finishes setting up stuff after the ui has been built for the
// first time.
func (t *tui) postCreate(g *gotui.Gui) error {
	if t.client.Config.Client.Syslog.LogLevel != "TRACE" {
		log.Tracef("setting up error listener")
		if err := loggo.RegisterWriter("tui", loggo.NewMinimumLevelWriter(loggocolor.NewColorWriter(t.errs), loggo.WARNING)); err == nil {
			t.errs.AddPostWriteHook(func(line *HistoryLine) error {
				t.createModal("Stimmtausch error", line.Text)
				return nil
			})
		} else {
			log.Errorf("error setting up error handler: %v", err)
		}
	} else {
		log.Tracef("not setting up error listener because I expect you're watching logs")
	}

	log.Tracef("setting up sent buffer to write to active connection")
	t.sent.AddPostWriteHook(func(line *HistoryLine) error {
		if t.currView == nil || !t.currView.connected {
			// The only case in which the UI dispatches is if there's no client
			// to do so. We want the client to do it usually this UI may not be
			// the only one.
			if line.Text[0] == '/' {
				s := strings.SplitN(line.Text[1:], " ", 2)
				if len(s) == 1 {
					s = append(s, "")
				}
				go t.client.Env.Dispatch(s[0], s[1])
				return nil
			}
			log.Warningf("no current connection!")
			return nil
		}
		if !t.currView.connected {
			return nil
		}
		_, err := fmt.Fprintln(t.currView.conn, line.Text)
		if err != nil {
			log.Warningf("error writing to connection")
			return errgo.Mask(err)
		}
		return nil
	})

	t.ready <- true

	return nil
}

// layout acts as the layout manager for the gotui.Gui, creating views.
func (t *tui) layout(g *gotui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("send", 0, maxY-5, maxX-1, maxY-1); err != nil {
		if err != gotui.ErrUnknownView {
			log.Warningf("unable to create view %+v", err)
			return errgo.Mask(err)
		} else {
			for i, view := range t.views {
				if err := view.updateRecvOrigin(i, g, t); err != nil {
					return errgo.Mask(err)
				}
			}
		}
		v.Editable = true
		v.Wrap = true
		if _, err := g.SetCurrentView("send"); err != nil {
			return errgo.Mask(err)
		}
		oldEditor := v.Editor
		v.Editor = gotui.EditorFunc(func(v *gotui.View, key gotui.Key, ch rune, mod gotui.Modifier) {
			oldEditor.Edit(v, key, ch, mod)
			go g.Update(func(gg *gotui.Gui) error {
				return t.updateCharCount(gg, v)
			})
		})
	}
	if v, err := g.SetView("charcount", maxX-16, maxY-2, maxX-2, maxY); err != nil {
		if err != gotui.ErrUnknownView {
			log.Warningf("unable to create view %+v", err)
			return errgo.Mask(err)
		}
		v.Frame = false
		fmt.Fprint(v, strings.Repeat("─", 10)+" 0 ")
	}
	if v, err := g.SetView("title", 2, maxY-6, t.titleLen+3, maxY-4); err != nil {
		if err != gotui.ErrUnknownView {
			log.Warningf("unable to create view %+v", err)
			return errgo.Mask(err)
		}
		v.Frame = false
		fmt.Fprint(v, " No world ")
		g.Update(t.postCreate)
	}
	return nil
}

// onResize performs operations that need to be accomplished when the window is
// resized. For now, it is very simple and just redraws the whole screen, but
// should, in the future, perform more expensive wrapping functions.
func (t *tui) onResize(g *gotui.Gui, x, y int) error {
	if t.currView == nil {
		return nil
	}
	v, err := g.View(t.currView.viewName)
	if err != nil {
		return errgo.Mask(err)
	}
	return errgo.Mask(t.redraw(g, v))
}

// updateSendTitle updates the title of the input buffer frame to show the
// world list with the active world and inactive worlds specified differently.
func (t *tui) updateSendTitle() {
	_, maxY := t.g.Size()
	if len(t.views) == 0 {
		t.title = "No world"
		t.titleLen = len(t.title) + 2
	} else {
		conns := make([]string, len(t.views))
		sep := " | "
		t.titleLen = 0 - len(sep) + 2
		for i, v := range t.views {
			title := v.displayName
			if v.hasMore {
				title = fmt.Sprintf("%s (+%d)", v.displayName, v.more)
			}
			t.titleLen += len(title) + len(sep)
			c, ok := t.client.Conn(v.connName)
			if !ok {
				continue
			}
			connected := c.Connected
			v.connected = connected
			if v.current {
				if connected {
					if v.hasMore {
						conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.ActiveMore, title)
					} else {
						conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.Active, title)
					}
				} else {
					if v.hasMore {
						conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.DisconnectedMore, title)
					} else {
						conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.DisconnectedActive, title)
					}
				}
			} else {
				if connected {
					if v.hasMore {
						conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.InactiveMore, title)
					} else {
						conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.Inactive, title)
					}
				} else {
					if v.hasMore {
						conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.DisconnectedMore, title)
					} else {
						conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.Disconnected, title)
					}
				}
			}
		}
		t.title = strings.Join(conns, sep)
	}
	log.Tracef("setting title to %s", t.title)
	if v, err := t.g.SetView("title", 1, maxY-6, t.titleLen+3, maxY-4); err == nil {
		t.g.Update(func(_ *gotui.Gui) error {
			v.Clear()
			fmt.Fprintf(v, " %s ", t.title)
			return nil
		})
	} else {
		log.Errorf("error setting title %v", err)
	}
}

// switchConn switches to a different connection, either by rotating the
// stack of connections or by switching to the given name.
func (t *tui) switchConn(action, conn string) error {
	if len(t.views) == 1 {
		return nil
	}
	if action == "active" {
		log.Tracef("rotating to active conn %s", conn)
		if conn == "1" {
			for i := t.currViewIndex + 1; i != t.currViewIndex; i++ {
				if i >= len(t.views) {
					i = -1
					continue
				}
				if t.views[i].hasMore {
					return errgo.Mask(t.switchConn("switch", t.views[i].connName))
				}
			}
		} else if conn == "-1" {
			for i := t.currViewIndex - 1; i != t.currViewIndex; i-- {
				if i < 0 {
					i = len(t.views)
					continue
				}
				if t.views[i].hasMore {
					return errgo.Mask(t.switchConn("switch", t.views[i].connName))
				}
			}
		}
	} else if action == "rotate" {
		i, err := strconv.Atoi(conn)
		if err != nil {
			return errgo.Mask(err)
		}
		log.Tracef("rotating conns %d", i)
		t.views[t.currViewIndex].current = false
		t.currViewIndex += i
		if t.currViewIndex < 0 {
			t.currViewIndex = len(t.views) + t.currViewIndex
		} else if t.currViewIndex >= len(t.views) {
			if len(t.views) > 1 {
				t.currViewIndex %= len(t.views)
			} else {
				t.currViewIndex = 0
			}
		}
		t.views[t.currViewIndex].current = true
	} else if action == "switch" {
		found := false
		var newIndex int
		for i, v := range t.views {
			if v.connName == conn {
				log.Tracef("switching to %s", conn)
				v.current = true
				newIndex = i
				found = true
				break
			}
		}
		if found {
			t.views[t.currViewIndex].current = false
			t.currViewIndex = newIndex
		} else {
			return errgo.Newf("connection %s not found", conn)
		}
	} else {
		return errgo.Newf("received unexpected action %s", action)
	}
	t.currView = t.views[t.currViewIndex]
	for _, v := range t.views {
		if err := v.updateRecvOrigin(t.currViewIndex, t.g, t); err != nil {
			return errgo.Mask(err)
		}
	}
	t.updateSendTitle()
	return nil
}

// removeWorld removes a world from the UI.
func (t *tui) removeWorld(world string) error {
	var r *receivedView
	for i, v := range t.views {
		if v.connName == world {
			r = v
			t.views = append(t.views[:i], t.views[i+1:]...)
			i++
			if len(t.views) > 0 {
				t.switchConn("rotate", "1")
			} else {
				t.currView = nil
			}
			break
		}
	}
	if r != nil {
		log.Infof("attempting to removew view %s", r.viewName)
		go t.g.Update(func(g *gotui.Gui) error {
			if err := g.DeleteView(r.viewName); err != nil {
				log.Errorf("unable to delete view %s: %v", r.viewName, err)
			}
			t.updateSendTitle()
			return nil
		})
	}
	return nil
}

// listen listens for events from the signal environment, then does nothing (but
// does it splendidly)
func (t *tui) listen() {
	for {
		res := <-t.listener
		switch res.Name {
		case "fg", ">", "<", "[", "]":
			// Switch to active view or error if not found.
			if err := t.switchConn(res.Payload[0], res.Payload[1]); err != nil {
				log.Warningf("error switching connections: %v", err)
			}
		case "connect", "c":
			// Maybe create receivedView here á là TF's trying to connect
			// screen? Maybe not.
			continue
		case "disconnect", "dc":
			// If it's a disconnect without a payload, redispatch with the
			// current connection's name.
			if (len(res.Payload) == 1 && res.Payload[0] != "-r") || len(res.Payload) != 0 || t.currView == nil {
				continue
			}
			res.Payload = append(res.Payload, t.currView.connName)
			log.Tracef("disconnecting current world %+v", res)
			go t.client.Env.DirectDispatch(res)
		case "help":
			// get the command text and tell the system to display it in a modal
			var cmd string
			if len(res.Payload) == 0 {
				cmd = "help"
			} else {
				cmd = res.Payload[0]
			}
			h, ok := help.HelpMessages[cmd]
			if !ok {
				log.Warningf("no help available for /%s", cmd)
				continue
			}
			helpText := help.RenderText(h)
			go t.client.Env.Dispatch("_client:showModal", fmt.Sprintf("Help: %s::\n%s", cmd, helpText))
		case "_client:connect":
			// Attach receivedView to conn.
			err := t.connect(res.Payload[0], t.g)
			if err != nil {
				log.Errorf("error setting up connection in ui: %v", err)
			}
		case "_client:disconnect":
			// Grey out tab in send title, grey out text in receivedView.
			t.updateSendTitle()
		case "_client:disconnected":
			// Grey out tab in send title, grey out text in receivedView.
			t.updateSendTitle()
		case "_client:allDisconnect":
			// do we really need to do anything?
			t.updateSendTitle()
		case "_client:quitReady":
			go t.g.Update(func(g *gotui.Gui) error {
				return gotui.ErrQuit
			})
		case "_client:showModal":
			res.Name = "_tui:showModal"
			go t.client.Env.DirectDispatch(res)
		case "_tui:showModal":
			t.createModal(res.Payload[0], res.Payload[1])
		case "_client:removeWorld", "remove", "r":
			if len(res.Payload) != 1 {
				log.Warningf("tried to remove a world without an argument")
				continue
			}
			world := res.Payload[0]
			if err := t.removeWorld(world); err != nil {
				log.Errorf("unable to remove world %s: %v", res.Payload[0], err)
			}
		default:
			log.Tracef("got unknown signal result %v", res)
			continue
		}
	}
}

func (t *tui) createModal(title, content string) {
	log.Tracef("showing modal with title %s", title)
	if t.modalOpen {
		return
	}
	go t.g.Update(func(g *gotui.Gui) error {
		maxX, maxY := g.Size()
		if v, err := g.SetView("modal", 3, 3, maxX-4, maxY-6); err != nil {
			if err != gotui.ErrUnknownView {
				log.Warningf("unable to create modal view %+v", err)
				return errgo.Mask(err)
			}
			t.modalOpen = true
			v.Frame = true
			v.FrameFgColor = gotui.ColorCyan | gotui.AttrBold
			v.Wrap = true
			v.WordWrap = true
			fmt.Fprint(v, content)
		}
		if v, err := g.SetView("modalTitle", 5, 2, len(title)+8, 4); err != nil {
			if err != gotui.ErrUnknownView {
				log.Warningf("unable to create modal view %+v", err)
				return errgo.Mask(err)
			}
			v.Frame = false
			fmt.Fprintf(v, " %s ", ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.ModalTitle, title))
		}
		modalHelpText := " Scroll: ↑/↓ | Close: <Enter> "
		if v, err := g.SetView("modalHelp", maxX-3-len(modalHelpText), maxY-7, maxX-6, maxY-5); err != nil {
			if err != gotui.ErrUnknownView {
				log.Warningf("unable to create modal view %+v", err)
				return errgo.Mask(err)
			}
			v.Frame = false
			fmt.Fprint(v, ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.ModalTitle, modalHelpText))
		}
		if _, err := g.SetCurrentView("modal"); err != nil {
			log.Warningf("unable to set current view %+v", err)
			return errgo.Mask(err)
		}
		return nil
	})
}

// GetMaxWidth returns the max requested width (or 100%).
func (t *tui) GetMaxWidth() int {
	maxX, _ := t.g.Size()
	maxWidth := t.client.Config.Client.UI.MaxWidth
	log.Debugf("maxX is %d and maxWidth is %d", maxX, maxWidth)
	if maxWidth <= 0 || maxWidth > maxX {
		maxWidth = maxX
	}
	return maxWidth
}

// Run creates the UI and starts the mainloop.
func (t *tui) Run(done, ready chan bool) {
	t.ready = ready

	log.Tracef("creating UI")
	var err error
	t.g, err = gotui.NewGui(gotui.Output256)
	if err != nil {
		log.Criticalf("unable to create ui: %v", err)
		os.Exit(2)
	}
	defer t.g.Close()

	t.g.Cursor = true
	t.g.Mouse = t.client.Config.Client.UI.Mouse

	t.g.SetManagerFunc(t.layout)
	t.g.SetResizeFunc(t.onResize)

	log.Tracef("adding keybindings")
	if err := t.keybindings(t.g); err != nil {
		log.Criticalf("ui couldn't create keybindings: %v", err)
		os.Exit(2)
	}

	log.Tracef("listening for signals")
	t.listener = make(chan signal.Signal)
	go t.listen()
	t.client.Env.AddListener("ui", t.listener)

	log.Tracef("running UI...")
	if err := t.g.MainLoop(); err != nil && err != gotui.ErrQuit {
		t.errs.Close()
		log.Criticalf("ui unexpectedly quit: %s: %s", err, errgo.Details(err))
		fmt.Printf("Oh no! Something went very wrong :( Here's all we know: %s: %s\n", err, errgo.Details(err))
		fmt.Println("Your connections will all close gracefully and any logs properly closed out.")
	}
	t.client.CloseAll()
	done <- true
}

// New instantiates a new Stimmtausch UI.
func New(c *client.Client) *tui {
	return &tui{
		client:   c,
		sent:     NewHistory(c.Config.Client.UI.History, false),
		errs:     NewHistory(100, true),
		title:    " No world ",
		titleLen: 10,
	}
}
