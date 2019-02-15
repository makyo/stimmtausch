// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package ui

import (
	"fmt"
	"os"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/makyo/gotui"

	"github.com/makyo/stimmtausch/client"
	"github.com/makyo/stimmtausch/macro"
)

type tui struct {
	g             *gotui.Gui
	client        *client.Client
	args          []string
	sent          *History
	currView      *receivedView
	currViewIndex int
	views         []*receivedView
	listener      chan macro.MacroResult
}

var log = loggo.GetLogger("stimmtausch.ui")

// connect tells the client to connect to the provided connection string, If
// successful, it will construct a receivedView to represent and hold that
// connection.
func (t *tui) connect(connectStr string, g *gotui.Gui) error {
	log.Tracef("creating a connection with connection string %s", connectStr)
	conn, err := t.client.Connect(connectStr)
	if err != nil {
		log.Errorf("unable to connect to %s: %v", connectStr, err)
		return err
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
			connName: conn.GetConnectionName(),
			conn:     conn,
			viewName: viewName,
			buffer:   NewHistory(t.client.Config.Client.UI.Scrollback),
			current:  true,
			index:    len(t.views),
		}
		t.views = append(t.views, t.currView)
		inputView, err := g.View("send")
		if err != nil {
			log.Warningf("unable to get send view to change title: %v", err)
		} else {
			inputView.Title = fmt.Sprintf(" %s ", conn.GetDisplayName())
		}

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

		log.Tracef("opening connection for %s", connectStr)
		err = conn.Open()
		if err != nil {
			log.Errorf("unable to open connection for %s: %v", connectStr, err)
			return err
		}
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
		_, err := fmt.Fprintln(t.currView.conn, line.Text)
		if err != nil {
			log.Warningf("error writing to connection")
			return err
		}
		return nil
	})

	// This is the real reason for this method. Connecting before everything
	// is set up means that some of the output from the worlds gets eaten by
	// the UI getting built.
	log.Tracef("attempting to connect with connection strings %v", t.args)
	for _, arg := range t.args {
		t.connect(arg, g)
	}
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
		v.Title = "No world"
		if _, err := g.SetCurrentView("send"); err != nil {
			return err
		}
	}
	return nil
}

// listen listens for events from the macro environment, then does nothing (but
// does it splendidly)
func (t *tui) listen() {
	for {
		<-t.listener
	}
}

func (t *tui) Run(done chan bool) {
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

	log.Tracef("listening for macros")
	t.listener = make(chan macro.MacroResult)
	t.client.Env.AddListener(t.listener)
	go t.listen()

	log.Tracef("running UI...")
	if err := t.g.MainLoop(); err != nil && err != gotui.ErrQuit {
		log.Criticalf("ui unexpectedly quit: %v", err)
	}
	t.client.CloseAll()
	done <- true
}

// New instantiates a new Stimmtausch UI.
func New(args []string, c *client.Client) *tui {
	return &tui{
		client: c,
		args:   args,
		sent:   NewHistory(c.Config.Client.UI.History),
	}
}
