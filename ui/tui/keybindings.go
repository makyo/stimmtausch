// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package tui

import (
	"fmt"
	"strings"

	"github.com/makyo/gotui"
)

// quit returns gotui.ErrQuit for when an event (such as the user quitting
// the programm) says the main loop should stop running.
func (t *tui) quit(g *gotui.Gui, v *gotui.View) error {
	return gotui.ErrQuit
}

// send sends whatever line is currently active in the input View to the
// sent buffer (and thus to the world via a post-write hook).
func (t *tui) send(g *gotui.Gui, v *gotui.View) error {
	buf := strings.TrimSpace(v.Buffer())
	if len(buf) == 0 {
		return nil
	}
	fmt.Fprint(t.sent, buf)
	v.Clear()
	v.SetCursor(0, 0)
	return nil
}

// arrowUp moves the cursor up in the input View. If there is text for the
// cursor to move up through, it will do so. If it's on the top line but not
// at the beginning of the line, it moves there. Otherwise, it attempts to
// scroll back through the sent history.
func (t *tui) arrowUp(g *gotui.Gui, v *gotui.View) error {
	cx, cy := v.Cursor()
	if cx == 0 {
		if cy == 0 {
			v.Clear()
			hl := t.sent.Back()
			if hl != nil {
				fmt.Fprint(v, hl.Text)
			}
		} else {
			v.SetCursor(0, 0)
		}
	} else {
		v.SetCursor(cx-1, cy)
	}
	return nil
}

// arrowDown moves the cursor down in the input View. If there is text for the
// cursor to move down through, it will do so. If it's on the last line but not
// at the end of the line, it moves there. Otherwise, it attempts to
// scroll forward through the sent history.
func (t *tui) arrowDown(g *gotui.Gui, v *gotui.View) error {
	cx, cy := v.Cursor()
	lines := v.ViewBufferLines()
	lineCount := len(v.ViewBufferLines()) - 1
	if lineCount == -1 {
		return nil
	}
	lastLineLen := len(lines[lineCount])
	if cx == lastLineLen || (cx == 0 && cy == 0 && !t.sent.OnLast()) {
		if cy == lineCount {
			v.Clear()
			v.SetCursor(0, 0)
			if !t.sent.OnLast() {
				fmt.Fprint(v, t.sent.Forward().Text)
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

// scrollConsole scrolls through text in the logging console.
// TODO this should be a receivedView rather than Autoscroll
func scrollConsole(v *gotui.View, delta int) {
	_, y := v.Origin()
	v.SetOrigin(0, y+delta)
}

// scrollUp scrolls the output buffer up by one screen. If that would go
// negative, it only scrolls to zero to prevent an error.
func (t *tui) scrollUp(g *gotui.Gui, v *gotui.View) error {
	v, err := g.View(t.currView.viewName)
	if err != nil {
		return err
	}
	_, y := v.Origin()
	_, maxY := v.Size()
	lines := len(v.ViewBufferLines())
	result := y - maxY
	if result < 0 {
		result = 0
	}
	log.Debugf("got %d lines, setting origin to 0,%d: %v", lines, result, v.SetOrigin(0, result))
	return nil
}

// scrollDown scrolls the output buffer down by one screen. If that would go
// past where the text is written, it scrolls by only that amount.
func (t *tui) scrollDown(g *gotui.Gui, v *gotui.View) error {
	v, err := g.View(t.currView.viewName)
	if err != nil {
		return err
	}
	_, y := v.Origin()
	_, maxY := v.Size()
	lines := len(v.ViewBufferLines())
	result := y + maxY
	if result < lines {
		if result+maxY > lines {
			result = result - (result + maxY - lines) - 1
		}
		if result != y {
			log.Debugf("got %d lines, setting origin to 0,%d: %v", lines, result, v.SetOrigin(0, result))
		}
	}
	return nil
}

// redraw forces a rerender of the current view in order to ensure that everything is in order
func (t *tui) redraw(g *gotui.Gui, v *gotui.View) error {
	log.Debugf("redrawing")
	v, err := g.View(t.currView.viewName)
	if err != nil {
		return err
	}
	x, y := v.Origin()
	v.Clear()
	fmt.Fprint(v, t.currView.buffer.String())
	// XXX This doesn't preserve, and I don't know why. Drat.
	// https://github.com/makyo/stimmtausch/issues/46
	v.SetOrigin(x, y)
	g.Update(func(gg *gotui.Gui) error {
		return t.currView.updateRecvOrigin(t.currViewIndex, gg)
	})
	return nil
}

// keybindings sets all keybindings used by the UI.
func (t *tui) keybindings(g *gotui.Gui) error {
	if err := g.SetKeybinding("", gotui.KeyCtrlC, gotui.ModNone, t.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyPgup, gotui.ModNone, t.scrollUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyPgdn, gotui.ModNone, t.scrollDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyCtrlL, gotui.ModNone, t.redraw); err != nil {
		return err
	}
	if err := g.SetKeybinding("send", gotui.KeyEnter, gotui.ModNone, t.send); err != nil {
		return err
	}
	if err := g.SetKeybinding("send", gotui.KeyArrowUp, gotui.ModNone, t.arrowUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("send", gotui.KeyArrowDown, gotui.ModNone, t.arrowDown); err != nil {
		return err
	}
	return nil
}