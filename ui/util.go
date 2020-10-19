package ui

import (
	"github.com/makyo/gotui"
)

// maybeSetView tries to set view positions intelligently. That is, if the coordinates
// provided are illegal, it will try to set them to what they already are.
func maybeSetView(g *gotui.Gui, title string, x1, y1, x2, y2 int) (*gotui.View, error) {
	if x2 <= x1+1 || y2 <= y1+1 {
		_x1, _y1, _x2, _y2, err := g.ViewPosition(title)
		if err == nil {
			x1 = _x1
			y1 = _y1
			x2 = _x2
			y2 = _y2
		}
	}
	return g.SetView(title, x1, y1, x2, y2)
}
