package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/juju/errgo"
	"github.com/makyo/ansigo"
	"github.com/makyo/gotui"
)

// send sends whatever line is currently active in the input View to the
// sent buffer (and thus to the world via a post-write hook).
func (t *tui) send(g *gotui.Gui, v *gotui.View) error {
	for i, l := range v.BufferLines() {
		buf := strings.TrimSpace(l)
		if i == 0 && len(buf) == 0 {
			// A single space is often used for defaulting values; allow that, but
			// otherwise trim the space. This is a hacky way to check, but it
			// appears that gotui won't fill the buffer with just spaces.
			if x, _ := v.Cursor(); x != 0 {
				buf = " "
			} else {
				return nil
			}
		}
		fmt.Fprint(t.sent, buf)
		time.Sleep(100 * time.Millisecond)
	}
	go g.Update(func(gg *gotui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		t.updateCharCount(gg, v)
		return nil
	})
	return nil
}

// forceNewline forces a newline in the send editor, letting you type more than
// one thing to send.
func (t *tui) forceNewline(g *gotui.Gui, v *gotui.View) error {
	v.EditNewLine()
	t.updateCharCount(g, v)
	return nil
}

// updateCharCount updates the character count
func (t *tui) updateCharCount(g *gotui.Gui, v *gotui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		line = ""
	}
	count := uint(len(strings.TrimSpace(line)))
	var max uint = 0
	if t.currView != nil && t.currView.conn != nil {
		max = t.currView.maxBuffer
	}
	countStr := fmt.Sprintf(" %d ", count)
	countLen := len(countStr)
	if max != 0 && count > max {
		countStr = ansigo.MaybeApplyOneWithReset("red", countStr)
	}
	countStr = strings.Repeat("─", 13-countLen) + countStr
	maxX, maxY := g.Size()
	go g.Update(func(gg *gotui.Gui) error {
		if vv, err := g.SetView("charcount", maxX-16, maxY-2, maxX-2, maxY); err != nil {
			return errgo.Mask(err)
		} else {
			vv.Clear()
			fmt.Fprint(vv, countStr)
		}
		return nil
	})
	return nil
}
