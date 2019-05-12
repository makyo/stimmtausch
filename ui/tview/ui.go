// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package tview

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
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
	screen tcell.Screen
	app    *tview.Application
	layout *tview.Grid
	views  *viewSet

	// The information for the title, which contains an appropriately colored list of worlds.
	title    string
	titleLen int
	titleBar *tview.TextView

	// The sent buffer and view.
	sent     *buffer.Buffer
	sentView *tview.Box

	// The client and listener used for interacting with the rest of the application.
	client   *client.Client
	listener chan signal.Signal
}

func (t *ui) update() {
	// update current view
	t.screen.Clear()
	t.layout.Clear().
		AddItem(tview.Primitive(t.views.currView.view), 0, 0, 1, 1, 0, 0, false).
		AddItem(tview.Primitive(t.titleBar), 1, 0, 1, 1, 0, 0, false).
		AddItem(tview.Primitive(t.sentView), 2, 0, 1, 1, 0, 0, false)
	// update send title
	_, _, width, _ := t.titleBar.GetInnerRect()
	t.titleBar.SetText(fmt.Sprintf("%s %s %s", string(tview.BoxDrawingsLightHorizontal), t.title, strings.Repeat(string(tview.BoxDrawingsLightHorizontal), width)))
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

	t.app = tview.NewApplication().
		SetScreen(t.screen)
	t.app.SetInputCapture(t.keybinding)
	t.layout = tview.NewGrid().
		SetRows(0, 1, 3).
		SetColumns(0).
		AddItem(tview.Primitive(t.titleBar), 1, 0, 1, 1, 0, 0, false).
		AddItem(tview.Primitive(t.sentView), 2, 0, 1, 1, 0, 0, false)

	log.Tracef("running UI...")
	if err := t.app.SetRoot(t.layout, true).Run(); err != nil {
		log.Criticalf("ui quit unexpectedly: %v", err)
	}
	t.client.CloseAll()
	done <- true
}

// New instantiates a new Stimmtausch UI.
func New(c *client.Client) (*ui, error) {
	tview.Styles.PrimitiveBackgroundColor = -1
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err = screen.Init(); err != nil {
		return nil, err
	}
	return &ui{
		client:   c,
		screen:   screen,
		sent:     buffer.New(c.Config.Client.UI.History),
		sentView: tview.NewBox(),
		views:    NewViewSet(),
		titleBar: tview.NewTextView().SetDynamicColors(true).SetRegions(true),
		title:    "No World",
		titleLen: 10,
	}, nil
}
