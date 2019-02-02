package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"

	"github.com/makyo/st/client"
)

type receivedView struct {
	connName string
	conn     io.WriteCloser
	viewName string
	buffer   *history
	current  bool
}

var (
	args     []string
	stClient *client.Client
	currView *receivedView
)

var (
	log           = loggo.GetLogger("stimmtausch.ui")
	sent          = NewHistory(1000)
	views         = []*receivedView{}
	currViewIndex = 0
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func send(g *gocui.Gui, v *gocui.View) error {
	buf := strings.TrimSpace(v.Buffer())
	if len(buf) == 0 {
		return nil
	}
	fmt.Fprintln(sent, buf)
	v.Clear()
	v.SetCursor(0, 0)
	return nil
}

func (v *receivedView) updateRecvSize(index int, g *gocui.Gui) error {
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
	g.Update(func(gg *gocui.Gui) error {
		if _, err := gg.SetView(v.viewName, recvX0-1, recvY0, recvX0+maxX, maxY-5); err != nil {
			log.Warningf("unable to create view %+v", err)
			return err
		}
		return nil
	})
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

func scrollConsole(v *gocui.View, delta int) {
	_, y := v.Origin()
	v.SetOrigin(0, y+delta)
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlLsqBracket, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		scrollConsole(v, -2)
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlRsqBracket, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		scrollConsole(v, 2)
		return nil
	}); err != nil {
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

func connect(connectStr string, g *gocui.Gui) error {
	log.Tracef("attempting to connect with connection string %s", connectStr)
	conn, err := stClient.Connect(connectStr)
	if err != nil {
		log.Errorf("unable to connect to %s: %v", connectStr, err)
	}
	viewName := fmt.Sprintf("recv%d", len(views))
	if v, err := g.SetView(viewName, -3, -3, -1, -1); err != nil {
		if err != gocui.ErrUnknownView {
			log.Warningf("unable to create view %+v", err)
			return err
		}
		fmt.Fprintln(v, "\na")
		v.Wrap = true
		v.Frame = false
		v.Autoscroll = true
		currView = &receivedView{
			connName: conn.GetConnectionName(),
			conn:     conn,
			viewName: viewName,
			buffer:   NewHistory(10000),
			current:  true,
		}
		views = append(views, currView)
		currViewIndex = len(views) - 1
		currView.buffer.AddPostWriteHook(func(line string) error {
			fmt.Fprint(v, line)
			return currView.updateRecvSize(currViewIndex, g)
		})
		conn.AddOutput(viewName, currView.buffer)
		err = conn.Open()
		if err != nil {
			return err
		}
	}
	return nil
}

func postCreate(g *gocui.Gui) error {
	console, err := g.View("console")
	if err != nil {
		return err
	}
	//loggo.RegisterWriter("console", loggocolor.NewWriter(console))
	loggo.ReplaceDefaultWriter(loggocolor.NewWriter(console))

	log.Tracef("setting up sent buffer to write to active connection")
	sent.AddPostWriteHook(func(line string) error {
		_, err := fmt.Fprintln(currView.conn, line)
		if err != nil {
			log.Warningf("error writing to connection")
			return err
		}
		return nil
	})
	log.Tracef("attempting to connect with connection strings %v", args)
	for _, arg := range args {
		connect(arg, g)
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("console", 0, 0, maxX-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			log.Warningf("unable to create view %+v", err)
			return err
		}
		v.Autoscroll = true
		postCreate(g)
	}
	if v, err := g.SetView("send", 0, maxY-5, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
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

func New(argsIn []string) {
	args = argsIn

	var err error
	stClient, err = client.New()
	if err != nil {
		log.Criticalf("could not create client: %v", err)
		os.Exit(4)
	}
	log.Tracef("created client: %+v", stClient)
	defer stClient.CloseAll()

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Criticalf("unable to create ui: %v", err)
		os.Exit(1)
	}
	defer g.Close()
	defer loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))

	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Criticalf("ui couldn't create keybindings: %v", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Criticalf("ui unexpectedly quit: %v", err)
		stClient.CloseAll()
		os.Exit(1)
	}
}
