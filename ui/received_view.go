package ui

import (
	"io"

	"github.com/juju/errgo"
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

	// Whether or not the connection is connected.
	connected bool

	// The connection's output buffer
	buffer *History

	// Whether or not the world is currently active.
	current bool

	// Whether or not we're in `more` mode
	hasMore bool

	// How many more lines we have
	more int

	// Our current origin Y
	scrollPos int

	// The index of the view.
	index int
}

// updateRecvOrigin updates the origin of every output gotui.View according to
// how many lines are in the buffer.
func (v *receivedView) updateRecvOrigin(index int, g *gotui.Gui, t *tui) error {
	maxX, maxY := g.Size()
	recvX0 := (maxX * index) - (maxX * v.index)
	if vv, err := g.SetView(v.viewName, recvX0-1, -1, recvX0+maxX, maxY-5); err != nil {
		log.Errorf("tried to set view to an invalid point (%d, %d) (%d %d)", recvX0-1, -1, recvX0+maxX, maxY-5)
		// return errgo.Mask(err)
		// Until https://github.com/makyo/stimmtausch/issues/115 is addressed, return nil.
		return nil
	} else {
		g.Update(func(gg *gotui.Gui) error {
			lines := len(vv.ViewBufferLines())
			_, vmaxY := vv.Size()
			_, y := vv.Origin()
			result := lines - vmaxY - 1
			if result != y && result >= 0 {
				if v.hasMore || !v.current {
					v.more = result - y
					v.hasMore = true
					log.Debugf("got more lines (now +%d) for %v", v.more, v.viewName)
				} else {
					v.more = 0
					v.hasMore = false
					if err = vv.SetOrigin(0, result); err != nil {
						return errgo.Mask(err)
					}
					log.Debugf("got %d lines, setting origin to 0,%d: %v", lines, result, err)
				}
				gg.Update(func(_ *gotui.Gui) error {
					t.updateSendTitle()
					return nil
				})
			}
			return nil
		})
	}
	return nil
}
