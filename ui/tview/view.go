package tview

import (
	"fmt"
	"io"

	"github.com/rivo/tview"

	"github.com/makyo/stimmtausch/buffer"
)

// view wraps a tview View and its corresponding buffer and connection.
type view struct {
	name        string
	displayName string
	conn        io.WriteCloser
	connected   bool
	writer      io.Writer
	view        *tview.TextView
	buffer      *buffer.Buffer
	current     bool
	app         *tview.Application

	// Mark whether the view is scrolled up. If so, add a marker where new
	// content starts.
	scrolled                 bool
	shouldAddScrollIndicator bool
}

// New creates a new view instance for the given connection, including its
// buffer and writer.
func NewView(application *tview.Application, name string, bufSize int) *view {
	tv := tview.NewTextView().
		SetWordWrap(true).
		SetScrollable(true).
		SetDynamicColors(true)
	w := tview.ANSIWriter(tv)

	buf := buffer.New(bufSize)
	buf.AddPostWriteHook(func(line *buffer.BufferLine) error {
		fmt.Fprint(w, line.Text)
		return nil
	})
	v := &view{
		name:   name,
		writer: w,
		view:   tv,
		buffer: buf,
		app:    application,
	}
	v.view.SetChangedFunc(v.maybeScrollToBottom)
	return v
}

// Write writes to the connection by way of the
func (v *view) Write(p []byte) (int, error) {
	return v.writer.Write(p)
}

func (v *view) maybeScrollToBottom() {
	if v.current && !v.scrolled {
		go v.app.QueueUpdateDraw(func() {
			v.view.ScrollToEnd()
		})
	} else {
		if v.shouldAddScrollIndicator {
			fmt.Fprint(v, "[-:-:bd]=====[-:-:-]")
			v.shouldAddScrollIndicator = false
		}
	}
}

func (v *view) scroll(amount int, byPage bool) {
	// Attempt to scroll by lines/page as requested in the given direction
	// (positive for down, negative for up).
	// If we're scrolling up and scrolled is false, set scrolled and
	// shouldAddScrollIndicator to true.
	// If we can't scroll down anymore, set scrolled to false.
}

func (v *view) redraw() {
	go v.app.QueueUpdateDraw(func() {
		v.view.Clear()
		go fmt.Fprint(v, v.buffer.String())
	})
}

// viewSet contains the collection of views associated with connections,
// whether active or inactive.
type viewSet struct {
	// The collection of views and their order.
	views map[string]*view
	order []string

	// The current view and its index in the order.
	currView  *view
	currIndex int
}

func NewViewSet() *viewSet {
	return &viewSet{
		views: map[string]*view{},
		order: []string{},
	}
}

// add adds a view to the viewSet and brings it to the foreground.
func (vs *viewSet) add(v *view) {
	vs.views[v.name] = v
	vs.order = append(vs.order, v.name)
	vs.fg(v.name)
}

// remove removes a view from the viewSet and attempts to cycle to the next
// available view.
func (vs *viewSet) remove(name string) error {
	if v, ok := vs.views[name]; ok {
		delete(vs.views, name)
		for i, n := range vs.order {
			if n == name {
				vs.order = append(vs.order[:i], vs.order[i+1:]...)
				break
			}
		}
		if vs.currView == v {
			vs.cycle(1)
		}
		return nil
	}
	return fmt.Errorf("no view named %s", name)
}

// cycle attemtps to cycle the requested amount of views - to the right for a
// positive number, and to the left for a negative one.
func (vs *viewSet) cycle(amount int) {
	vs.currIndex = amount % (len(vs.order) - 1)
	vs.fg(vs.order[vs.currIndex])
}

// fg attempts to bring the named view to the foreground. If there is no view
// by that name, it returns an error.
func (vs *viewSet) fg(name string) error {
	if v, ok := vs.views[name]; ok {
		v.current = true
		vs.currView = vs.views[name]
		for _, n := range vs.order {
			if n != name {
				vs.views[n].current = false
			}
		}
		return nil
	}
	return fmt.Errorf("no view named %s", name)
}
