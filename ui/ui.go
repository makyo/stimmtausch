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

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	ansi "github.com/makyo/ansigo"
	"github.com/makyo/gotui"

	"github.com/makyo/stimmtausch/client"
	"github.com/makyo/stimmtausch/help"
	"github.com/makyo/stimmtausch/signal"
)

type tui struct {
	g             *gotui.Gui
	client        *client.Client
	sent          *History
	currView      *receivedView
	currViewIndex int
	views         []*receivedView
	listener      chan signal.Signal
	ready         chan bool
	title         string
	titleLen      int
}

var log = loggo.GetLogger("stimmtausch.ui")

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
				return err
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
			return err
		}

		// Prime the view with a newline, which keeps it from complaining about
		// coordinates later.
		fmt.Fprintln(v, "\n ")
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
			buffer:      NewHistory(t.client.Config.Client.UI.Scrollback),
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
		t.currView.buffer.AddPostWriteHook(func(line *HistoryLine) error {
			fmt.Fprint(v, line.Text)
			g.Update(func(gg *gotui.Gui) error {
				return t.currView.updateRecvOrigin(t.currViewIndex, gg)
			})
			return nil
		})

		// Add the received history to the connection as an output.
		conn.AddOutput(viewName, t.currView.buffer, true)

		log.Tracef("opening connection for %s", name)
		err = conn.Open()
		if err != nil {
			log.Errorf("unable to open connection for %s: %v", name, err)
			return err
		}
		t.currView.connected = conn.Connected
		t.updateSendTitle()
	}
	return nil
}

// postCreate finishes setting up stuff after the ui has been built for the
// first time.
func (t *tui) postCreate(g *gotui.Gui) error {
	console, err := g.View("console")
	if err != nil {
		return err
	}

	loggo.RegisterWriter("console", loggocolor.NewColorWriter(console))

	log.Tracef("setting up sent buffer to write to active connection")
	t.sent.AddPostWriteHook(func(line *HistoryLine) error {
		if t.currView == nil || !t.currView.connected {
			// The only case in which the UI dispatches is if there's no client
			// to do so. We want the client to do it usually this UI may not be
			// the only one.
			if line.Text[0] == '/' {
				s := strings.SplitN(line.Text[1:], " ", 2)
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
			return err
		}
		return nil
	})

	t.ready <- true

	return nil
}

// layout acts as the layout manager for the gotui.Gui, creating views.
func (t *tui) layout(g *gotui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("console", 0, 0, maxX-1, 3); err != nil {
		if err != gotui.ErrUnknownView {
			log.Warningf("unable to create view %+v", err)
			return err
		}
		v.Autoscroll = true
		g.Update(t.postCreate)
	}
	if v, err := g.SetView("send", 0, maxY-5, maxX-1, maxY-1); err != nil {
		if err != gotui.ErrUnknownView {
			log.Warningf("unable to create view %+v", err)
			return err
		} else {
			for i, view := range t.views {
				if err := view.updateRecvOrigin(i, g); err != nil {
					return err
				}
			}
		}
		v.Editable = true
		v.Wrap = true
		if _, err := g.SetCurrentView("send"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("title", 2, maxY-6, t.titleLen+3, maxY-4); err != nil {
		if err != gotui.ErrUnknownView {
			log.Warningf("unable to create view %+v", err)
			return err
		}
		v.Frame = false
		fmt.Fprint(v, " No world ")
	}
	return nil
}

// updateSendTitle updates the title of the input buffer frame to show the
// world list with the active world and inactive worlds specified differently.
func (t *tui) updateSendTitle() {
	_, maxY := t.g.Size()
	conns := make([]string, len(t.views))
	sep := " | "
	t.titleLen = 0 - len(sep) + 2
	for i, v := range t.views {
		t.titleLen += len(v.displayName) + len(sep)
		c, ok := t.client.Conn(v.connName)
		if !ok {
			continue
		}
		connected := c.Connected
		v.connected = connected
		if v.current {
			if connected {
				conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.Active, v.displayName)
			} else {
				conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.DisconnectedActive, v.displayName)
			}
		} else {
			if connected {
				conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.Inactive, v.displayName)
			} else {
				conns[i] = ansi.MaybeApplyWithReset(t.client.Config.Client.UI.Colors.SendTitle.Disconnected, v.displayName)
			}
		}
	}
	t.title = strings.Join(conns, sep)
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
	if action == "rotate" {
		i, err := strconv.Atoi(conn)
		if err != nil {
			return err
		}
		log.Tracef("rotating conns %d", i)
		t.views[t.currViewIndex].current = false
		t.currViewIndex += i
		if t.currViewIndex < 0 {
			t.currViewIndex = len(t.views) + t.currViewIndex
		} else if t.currViewIndex >= len(t.views) {
			if len(t.views) > 1 {
				t.currViewIndex %= len(t.views) - 1
			} else {
				t.currViewIndex = 0
			}
		}
		t.views[t.currViewIndex].current = true
	} else if action == "switch" {
		if conn == "" {
			// TODO switch to the next active connection
		}
		found := false
		for i, v := range t.views {
			if v.connName == conn {
				log.Tracef("switchign to %s", conn)
				t.currViewIndex = i
				v.current = true
				found = true
				break
			}
		}
		if found {
			t.views[t.currViewIndex].current = false
		} else {
			return fmt.Errorf("connection %s not found", conn)
		}
	} else {
		return fmt.Errorf("received unexpected action %s", action)
	}
	t.currView = t.views[t.currViewIndex]
	for _, v := range t.views {
		if err := v.updateRecvOrigin(t.currViewIndex, t.g); err != nil {
			return err
		}
	}
	t.updateSendTitle()
	return nil
}

// listen listens for events from the signal environment, then does nothing (but
// does it splendidly)
func (t *tui) listen() {
	for {
		res := <-t.listener
		switch res.Name {
		case "fg", ">", "<":
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
			if len(res.Payload) != 0 {
				continue
			}
			res.Payload = []string{t.currView.connName}
			log.Tracef("disconnecting current world %+v", res)
			go t.client.Env.DirectDispatch(res)
		case "help":
			// get the command text and tell the system to display it in a modal
			cmd := res.Payload[0]
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
		case "_client:allDisconnect":
			// do we really need to do anything?
			t.updateSendTitle()
		case "_client:showModal":
			res.Name = "_tui:showModal"
			go t.client.Env.DirectDispatch(res)
		case "_tui:showModal":
			fmt.Fprintf(t.currView.buffer, "~~~~~ %s\n%s\n~~~~~\n", res.Payload[0], res.Payload[1])
		default:
			log.Tracef("got unknown signal result %v", res)
			continue
		}
	}
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
		log.Criticalf("ui unexpectedly quit: %v", err)
	}
	t.client.CloseAll()
	done <- true
}

// New instantiates a new Stimmtausch UI.
func New(c *client.Client) *tui {
	return &tui{
		client:   c,
		sent:     NewHistory(c.Config.Client.UI.History),
		title:    " No world ",
		titleLen: 10,
	}
}
