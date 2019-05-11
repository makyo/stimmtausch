// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package tview

import (
	"github.com/juju/loggo"
	"github.com/rivo/tview"

	"github.com/makyo/stimmtausch/buffer"
	"github.com/makyo/stimmtausch/client"
	"github.com/makyo/stimmtausch/signal"
)

var log = loggo.GetLogger("stimmtausch.ui.tview")

// ui contains all of the information necessary to run the Stimmtausch UI, from
// the client and listener to the tview application.
type ui struct {
	// The tview application and flex layout, along with the set of views.
	app   *tview.Application
	flex  *tview.Flex
	views *viewSet

	// The information for the title, which contains an appropriately colored list of worlds.
	title    string
	titleLen int

	// The sent buffer and view.
	sent     *buffer.Buffer
	sentView *tview.Box

	// The client and listener used for interacting with the rest of the application.
	client   *client.Client
	listener chan signal.Signal
}

func (t *ui) update() {
	// update current view
	// update send title
}

func (t *ui) connect(name string) error {
	// attach connection to view
	return nil
}

// Run creates the UI and starts the mainloop.
func (t *ui) Run(done chan bool) {
	t.listener = make(chan signal.Signal)
	go t.listen()
	t.client.Env.AddListener("ui", t.listener)

	t.app = tview.NewApplication()
	t.flex = tview.NewFlex()

	log.Tracef("running UI...")
	if err := t.app.SetRoot(t.flex, true).SetFocus(t.flex).Run(); err != nil {
		log.Criticalf("ui quit unexpectedly: %v", err)
	}
	t.client.CloseAll()
	done <- true
}

// New instantiates a new Stimmtausch UI.
func New(c *client.Client) *ui {
	return &ui{
		client:   c,
		sent:     buffer.New(c.Config.Client.UI.History),
		views:    &viewSet{},
		title:    " No World ",
		titleLen: 10,
	}
}
