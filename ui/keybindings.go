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
			fmt.Fprint(v, sent.back())
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

// scrollConsole scrolls through text in the logging console.
func scrollConsole(v *gotui.View, delta int) {
	_, y := v.Origin()
	v.SetOrigin(0, y+delta)
}

// keybindings sets all keybindings used by the UI.
func keybindings(g *gotui.Gui) error {
	if err := g.SetKeybinding("", gotui.KeyCtrlC, gotui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyCtrlLsqBracket, gotui.ModNone, func(g *gotui.Gui, v *gotui.View) error {
		scrollConsole(v, -2)
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyCtrlRsqBracket, gotui.ModNone, func(g *gotui.Gui, v *gotui.View) error {
		scrollConsole(v, 2)
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyCtrlL, gotui.ModNone, func(g *gotui.Gui, v *gotui.View) error {
		g.Update(func(g *gotui.Gui) error { return nil })
		log.Debugf("redrawing")
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
