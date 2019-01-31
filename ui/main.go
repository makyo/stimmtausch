package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"

	"github.com/makyo/st/client"
)

var log = loggo.GetLogger("stimmtausch.ui")
var sent = NewHistory(1000)
var received = NewHistory(10000)
var lines = 1

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func send(g *gocui.Gui, v *gocui.View) error {
	recv, err := g.View("recv")
	if err != nil {
		return err
	}
	buf := strings.TrimSpace(v.Buffer())
	if len(buf) == 0 {
		return nil
	}
	sent.add(buf)
	received.add(buf)
	v.Clear()
	v.SetCursor(0, 0)
	fmt.Fprintf(recv, "\n%s", received.current())
	g.Update(updateRecvSize)
	return nil
}

func updateRecvSize(g *gocui.Gui) error {
	v, err := g.View("recv")
	if err != nil {
		return err
	}
	lines = len(v.ViewBufferLines())
	return nil
}

func arrowUp(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	if cx == 0 {
		if cy == 0 {
			v.Clear()
			fmt.Fprint(v, sent.back())
		} else {
			v.SetCursor(0, 0)
		}
	} else {
		v.SetCursor(cx-1, cy)
	}
	return nil
}

func arrowDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	lines := v.ViewBufferLines()
	lineCount := len(v.ViewBufferLines()) - 1
	if lineCount == -1 {
		return nil
	}
	lastLineLen := len(lines[lineCount])
	if cx == lastLineLen || (cx == 0 && cy == 0 && !sent.onLast()) {
		if cy == lineCount {
			v.Clear()
			v.SetCursor(0, 0)
			if !sent.onLast() {
				fmt.Fprint(v, sent.forward())
			}
		} else {
			if lineCount == 0 {
				v.SetCursor(cx, cy+1)
			}
		}
	} else {
		if cy == lineCount {
			v.SetCursor(lastLineLen, cy)
		} else {
			v.SetCursor(cx, cy+1)
		}
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("send", gocui.KeyEnter, gocui.ModNone, send); err != nil {
		return err
	}
	if err := g.SetKeybinding("send", gocui.KeyArrowUp, gocui.ModNone, arrowUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("send", gocui.KeyArrowDown, gocui.ModNone, arrowDown); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	recvY0 := 3
	if lines < maxY-7 {
		recvY0 = maxY - 5 - lines
	}
	if v, err := g.SetView("recv", -1, recvY0, maxX, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, "\n")
		v.Wrap = true
		v.Frame = false
		v.Autoscroll = true
	}
	if v, err := g.SetView("console", 0, 0, maxX-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		loggo.RegisterWriter("ui view", loggocolor.NewWriter(v))
	}
	if v, err := g.SetView("send", 0, maxY-5, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
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

func New(args []string) {
	c, err := client.New()
	if err != nil {
		log.Criticalf("could not create client: %v", err)
		os.Exit(4)
	}
	log.Tracef("created client: %+v", c)

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Criticalf("unable to create ui: %v", err)
		os.Exit(1)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Criticalf("ui couldn't create keybindings: %v", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Criticalf("ui unexpectedly quit: %v", err)
		os.Exit(1)
	}
}
