// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package ui

import (
	"fmt"

	"github.com/makyo/gotui"
)

// quit returns gotui.ErrQuit for when an event (such as the user quitting
// the programm) says the main loop should stop running.
func quit(g *gotui.Gui, v *gotui.View) error {
	return gotui.ErrQuit
}

// arrowUp moves the cursor up in the input View. If there is text for the
// cursor to move up through, it will do so. If it's on the top line but not
// at the beginning of the line, it moves there. Otherwise, it attempts to
// scroll back through the sent history.
func arrowUp(g *gotui.Gui, v *gotui.View) error {
	cx, cy := v.Cursor()
	if cx == 0 {
		if cy == 0 {
			v.Clear()
			hl := sent.back()
			if hl != nil {
				fmt.Fprint(v, hl.text)
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
func arrowDown(g *gotui.Gui, v *gotui.View) error {
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
				fmt.Fprint(v, sent.forward().text)
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
func scrollConsole(v *gotui.View, delta int) {
	_, y := v.Origin()
	v.SetOrigin(0, y+delta)
}

func scrollUp(g *gotui.Gui, v *gotui.View) error {
	v, err := g.View(currView.viewName)
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

func scrollDown(g *gotui.Gui, v *gotui.View) error {
	v, err := g.View(currView.viewName)
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

// keybindings sets all keybindings used by the UI.
func keybindings(g *gotui.Gui) error {
	if err := g.SetKeybinding("", gotui.KeyCtrlC, gotui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyPgup, gotui.ModNone, scrollUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyPgdn, gotui.ModNone, scrollDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyCtrlL, gotui.ModNone, func(g *gotui.Gui, v *gotui.View) error {
		log.Debugf("redrawing")
		v, err := g.View(currView.viewName)
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprint(v, currView.buffer.String())
		g.Update(func(gg *gotui.Gui) error {
			return currView.updateRecvOrigin(currViewIndex, gg)
		})
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("send", gotui.KeyEnter, gotui.ModNone, send); err != nil {
		return err
	}
	if err := g.SetKeybinding("send", gotui.KeyArrowUp, gotui.ModNone, arrowUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("send", gotui.KeyArrowDown, gotui.ModNone, arrowDown); err != nil {
		return err
	}
	return nil
}
