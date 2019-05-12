package tview

import (
	"github.com/gdamore/tcell"
)

func (t *ui) keybinding(evt *tcell.EventKey) *tcell.EventKey {
	if evt.Key() == tcell.KeyCtrlL {
		go t.app.QueueUpdateDraw(func() {
			t.views.currView.redraw()
			t.screen.Clear()
			go t.app.QueueUpdateDraw(t.update)
		})
	}
	return evt
}
