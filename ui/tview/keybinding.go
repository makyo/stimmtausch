package tview

import (
	"github.com/gdamore/tcell"
)

func (t *ui) keybinding(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlL:
		go t.app.QueueUpdateDraw(func() {
			t.views.currView.redraw()
			t.screen.Clear()
			go t.app.QueueUpdateDraw(t.update)
		})
		return nil
	case tcell.KeyPgUp, tcell.KeyPgDn:
		return t.scroll(event)
	case tcell.KeyEnter:
		return t.maybeSend(event)
	}
	return event
}

func (t *ui) scroll(event *tcell.EventKey) *tcell.EventKey {
	return nil
}

func (t *ui) maybeSend(event *tcell.EventKey) *tcell.EventKey {
	if event.Modifiers() != tcell.ModNone {
		if event.Modifiers() == tcell.ModCtrl {
			return event
		}
		return nil
	}
	// Send...
	return nil
}
