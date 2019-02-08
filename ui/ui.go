// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/makyo/gotui"

	"github.com/makyo/stimmtausch/client"
	"github.com/makyo/stimmtausch/config"
)

// recievedView represents the gotui view which holds text received from the
// connection (the output buffer, basically).
type receivedView struct {
	// The name of the connection
	connName string

	// The name of the gotui View
	viewName string

	// The connection object itself in the form of an io.WriteCloser.
	conn io.WriteCloser

	// The connection's output buffer
	buffer *history

	// Whether or not the world is currently active.
	current bool
}

var (
	args     []string
	stClient *client.Client
	currView *receivedView
	sent     *history
	cfg      *config.Config
)

var (
	log           = loggo.GetLogger("stimmtausch.ui")
	views         = []*receivedView{}
	currViewIndex = 0
)

// send sends whatever line is currently active in the input View to the
// sent buffer (and thus to the world via a post-write hook).
func send(g *gotui.Gui, v *gotui.View) error {
	buf := strings.TrimSpace(v.Buffer())
	if len(buf) == 0 {
		return nil
	}
	fmt.Fprint(sent, buf)
	v.Clear()
	v.SetCursor(0, 0)
	return nil
}

// updateRecvSize updates the size of every output gotui.View according to
// how many lines are in the buffer. This is how we mock having the text
// scroll up from the bottom before the buffer gets to be larger than the
// window size, as text is always written from the top of the view down.
func (v *receivedView) updateRecvSize(index int, g *gotui.Gui) error {
	view, err := g.View(v.viewName)
	if err != nil {
		return err
	}
	lines := len(view.ViewBufferLines())
	maxX, maxY := g.Size()
	recvY0 := 3
	if lines < maxY-7 {
		recvY0 = maxY - 7 - lines
	}
	recvX0 := (maxX * index) - (maxX * currViewIndex)
	g.Update(func(gg *gotui.Gui) error {
		if _, err := gg.SetView(v.viewName, recvX0-1, recvY0, recvX0+maxX, maxY-5); err != nil {
			log.Warningf("unable to create view %+v", err)
			return err
		}
		return nil
	})
	return nil
}

// connect tells the client to connect to the provided connection string, If
// successful, it will construct a receivedView to represent and hold that
// connection.
func connect(connectStr string, g *gotui.Gui) error {
	log.Tracef("creating a connection with connection string %s", connectStr)
	conn, err := stClient.Connect(connectStr)
	if err != nil {
		log.Errorf("unable to connect to %s: %v", connectStr, err)
		return err
	}

	viewName := fmt.Sprintf("recv%d", len(views))
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
		v.IndentFirst = cfg.Client.UI.IndentFirst
		v.IndentSubsequent = cfg.Client.UI.IndentSubsequent
		v.Frame = false
		v.Autoscroll = true
		currView = &receivedView{
			connName: conn.GetConnectionName(),
			conn:     conn,
			viewName: viewName,
			buffer:   NewHistory(cfg.Client.UI.Scrollback),
			current:  true,
		}
		views = append(views, currView)

		// Set this view's index as the current view index.
		currViewIndex = len(views) - 1

		// Attach a hook that writes to the view when a line is received in
		// the received history for the connection.
		currView.buffer.AddPostWriteHook(func(line *historyLine) error {
			fmt.Fprint(v, line.text)
			return currView.updateRecvSize(currViewIndex, g)
		})

		// Add the received history to the connection as an output.
		conn.AddOutput(viewName, currView.buffer)

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
func postCreate(g *gotui.Gui) error {
	console, err := g.View("console")
	if err != nil {
		return err
	}

	// XXX Leaving the default writer in place writes to stderr, which
	// obviously writes to the screen. However, replacing it means that a lot
	// of the logging disappears. Currently, running with, e.g:
	//     go run main.go world 2>log.out
	// and then tailing the log in a separate window works, but bleh. Need to
	// either figure out logging to a file or why stderr shows through the UI.
	// Probably both.
	loggo.RegisterWriter("console", loggocolor.NewWriter(console))
	//loggo.ReplaceDefaultWriter(loggocolor.NewWriter(console))

	log.Tracef("setting up sent buffer to write to active connection")
	sent.AddPostWriteHook(func(line *historyLine) error {
		_, err := fmt.Fprintln(currView.conn, line.text)
		if err != nil {
			log.Warningf("error writing to connection")
			return err
		}
		return nil
	})

	// This is the real reason for this method. Connecting before everything
	// is set up means that some of the output from the worlds gets eaten by
	// the UI getting built.
	log.Tracef("attempting to connect with connection strings %v", args)
	for _, arg := range args {
		connect(arg, g)
	}
	return nil
}

// layout acts as the layout manager for the gotui.Gui, creating views.
func layout(g *gotui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("console", 0, 0, maxX-1, 3); err != nil {
		if err != gotui.ErrUnknownView {
			log.Warningf("unable to create view %+v", err)
			return err
		}
		v.Autoscroll = true
		g.Update(postCreate)
	}
	if v, err := g.SetView("send", 0, maxY-5, maxX-1, maxY-1); err != nil {
		if err != gotui.ErrUnknownView {
			log.Warningf("unable to create view %+v", err)
			return err
		} else {
			for i, view := range views {
				if err := view.updateRecvSize(i, g); err != nil {
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

// New instantiates a new Stimmtausch UI.
func New(_args []string, _cfg *config.Config) {
	args = _args
	cfg = _cfg

	sent = NewHistory(cfg.Client.UI.History)

	log.Tracef("creating client")
	var err error
	stClient, err = client.New(cfg)
	if err != nil {
		log.Criticalf("could not create client: %v", err)
		os.Exit(2)
	}
	log.Tracef("created client: %+v", stClient)
	defer stClient.CloseAll()

	log.Tracef("creating UI")
	g, err := gotui.NewGui(gotui.Output256)
	if err != nil {
		log.Criticalf("unable to create ui: %v", err)
		os.Exit(2)
	}
	defer g.Close()

	// XXX See above note on default writers.
	defer loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))

	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(layout)

	log.Tracef("adding keybindings")
	if err := keybindings(g); err != nil {
		log.Criticalf("ui couldn't create keybindings: %v", err)
	}

	log.Tracef("running UI...")
	if err := g.MainLoop(); err != nil && err != gotui.ErrQuit {
		log.Criticalf("ui unexpectedly quit: %v", err)
		stClient.CloseAll()
		os.Exit(2)
	}
}
