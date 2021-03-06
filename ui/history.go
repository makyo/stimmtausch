// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package ui

import (
	"strings"
	"time"
)

// HistoryLine represents a timestamped line of text.
type HistoryLine struct {
	Timestamp time.Time
	Text      string
}

// History represents a rolling buffer of lines used for input and output.
type History struct {
	// The index of the current line
	curr int

	// The maximum number of lines to hold in the buffer.
	max int

	// The lines we're keeping track of.
	lines []*HistoryLine

	// A list of functions to execute whenever the buffer is written to.
	postWriteHooks []func(*HistoryLine) error

	// Whether or not the buffer is closed.
	closed bool

	// Whether or not the current line is complete.
	lineComplete bool

	// Whether or not the writer should only write on newlines
	allowFragments bool
}

// add appends a line to the history, rolling a line out if necessary.
func (h *History) add(line string) bool {
	if h.allowFragments {
		containsNewline := strings.Contains(line, "\n")
		if !h.lineComplete {
			h.appendToCurrent(line)
			if containsNewline {
				h.lineComplete = true
				return true
			} else {
				return false
			}
		}
		h.lineComplete = containsNewline
	}
	l := &HistoryLine{
		Timestamp: time.Now(),
		Text:      line,
	}
	h.lines = append(h.lines, l)
	if len(h.lines) > h.max {
		h.lines = h.lines[1 : h.max+1]
	}
	h.curr = len(h.lines) - 1
	if h.allowFragments {
		return h.lineComplete
	}
	return true
}

// appendToCurrent appends a string to the current line.
func (h *History) appendToCurrent(line string) {
	h.lines[h.curr].Text = h.lines[h.curr].Text + line
}

// Current returns the current line in the buffer.
func (h *History) Current() *HistoryLine {
	if len(h.lines) == 0 {
		return nil
	}
	return h.lines[h.curr]
}

// Forward moves the cursor forward in time one line and returns the current
// line's content.
func (h *History) Forward() *HistoryLine {
	h.curr++
	if h.curr >= len(h.lines) {
		h.curr = len(h.lines) - 1
	}
	return h.Current()
}

// Back moves the cursor back in time one line and returns the current
// line's content.
func (h *History) Back() *HistoryLine {
	if h.curr < 0 {
		h.curr = 0
	}
	line := h.Current()
	h.curr--
	return line
}

// Last moves the cursor to the last line.
func (h *History) Last() *HistoryLine {
	h.curr = len(h.lines) - 1
	return h.Current()
}

// onLast returns true if the current line is the last (most recent) line.
func (h *History) onLast() bool {
	return h.curr == len(h.lines)-1
}

// String outputs the entire buffer as it stands.
func (h *History) String() string {
	var b strings.Builder
	for _, l := range h.lines {
		b.WriteString(l.Text)
	}
	return b.String()
}

// Write accepts a byte array and appends it to the buffer. It then executes
// every post-write hook. It returns the number of bytes written and any
// errors that occured.
// Fulfills io.Writer
func (h *History) Write(line []byte) (int, error) {
	if h.closed {
		return -1, nil
	}
	log.Tracef("received %d bytes", len(line))
	if h.add(string(line)) {

		// This currently happens synchronously and serially. Should make sure we
		// want to do it this way or use a context.
		for _, hook := range h.postWriteHooks {
			err := hook(h.Current())
			if err != nil {
				return len(line), err
			}
		}
	}
	return len(line), nil
}

func (h *History) Size() int {
	return len(h.lines)
}

// Close does nothing, and does it splendidly.
// Fulfills io.Closer
func (h *History) Close() error {
	h.closed = true
	return nil
}

// AddPostWriteHook accepts a function to be run whenever a write to the buffer
// succeeds.
func (h *History) AddPostWriteHook(f func(*HistoryLine) error) {
	if h.closed {
		return
	}
	h.postWriteHooks = append(h.postWriteHooks, f)
}

// NewHistory returns a new history buffer.
func NewHistory(max int, allowFragments bool) *History {
	return &History{
		curr:           0,
		max:            max,
		lines:          []*HistoryLine{},
		lineComplete:   true,
		allowFragments: allowFragments,
	}
}
