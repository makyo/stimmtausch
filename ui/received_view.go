package ui

import (
	"io"

	"github.com/makyo/gotui"
)

// recievedView represents the gotui view which holds text received from the
// connection (the output buffer, basically).
type receivedView struct {
	// The name of the connection
	connName string

	// The display name to show in the UI
	displayName string

	// The name of the gotui View
	viewName string

	// The connection object itself in the form of an io.WriteCloser.
	conn io.WriteCloser

	// The connection's output buffer
	buffer *History

	// Whether or not the world is currently active.
	current bool

	// The index of the view.
	index int
}

// updateRecvOrigin updates the origin of every output gotui.View according to
// how many lines are in the buffer.
func (v *receivedView) updateRecvOrigin(index int, g *gotui.Gui) error {
	maxX, maxY := g.Size()
	recvX0 := (maxX * index) - (maxX * v.index)
	if vv, err := g.SetView(v.viewName, recvX0-1, 3, recvX0+maxX, maxY-5); err != nil {
		return err
	} else {
		g.Update(func(gg *gotui.Gui) error {
			lines := len(vv.ViewBufferLines())
			_, vmaxY := vv.Size()
			_, y := vv.Origin()
			result := lines - vmaxY - 1
			if result != y && result >= 0 {
				log.Debugf("got %d lines, setting origin to 0,%d: %v", lines, result, vv.SetOrigin(0, result))
			}
			return nil
		})
	}
	return nil
}
