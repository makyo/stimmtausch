// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package ui

import (
	"strings"
)

// history represents a rolling buffer of lines used for input and output.
type history struct {
	// The index of the current line
	curr int

	// The maximum number of lines to hold in the buffer.
	max int

	// The lines we're keeping track of.
	lines []string

	// A list of functions to execute whenever the buffer is written to.
	postWriteHooks []func(string) error
}

// add appends a line to the history, rolling a line out if necessary.
func (h *history) add(line string) {
	h.lines = append(h.lines, strings.TrimSpace(line))
	if len(h.lines) > h.max {
		h.lines = h.lines[:h.max]
	}
	h.curr = len(h.lines) - 1
}

// current returns the current line in the buffer.
func (h *history) current() string {
	if len(h.lines) == 0 {
		return ""
	}
	return h.lines[h.curr]
}

// forward moves the cursor forward in time one line and returns the current
// line's content.
func (h *history) forward() string {
	h.curr++
	if h.curr >= len(h.lines) {
		h.curr = len(h.lines) - 1
	}
	return h.current()
}

// back moves the cursor back in time one line and returns the current
// line's content.
func (h *history) back() string {
	if h.curr < 0 {
		h.curr = 0
	}
	line := h.current()
	h.curr--
	return line
}

// onLast returns true if the current line is the last (most recent) line.
func (h *history) onLast() bool {
	return h.curr == len(h.lines)-1
}

// String outputs the entire buffer as it stands.
func (h *history) String() string {
	return strings.TrimSpace(strings.Join(h.lines, ""))
}

// Write accepts a byte array and appends it to the buffer. It then executes
// every post-write hook. It returns the number of bytes written and any
// errors that occured.
// Fulfills io.Writer
func (h *history) Write(line []byte) (int, error) {
	log.Tracef("received %d bytes", len(line))
	h.add(string(line))

	// This currently happens synchronously and serially. Should make sure we
	// want to do it this way or use a context.
	for _, hook := range h.postWriteHooks {
		err := hook(string(line))
		if err != nil {
			return len(line), err
		}
	}
	return len(line), nil
}

// Close does nothing, and does it splendidly.
// Fulfills io.Closer
func (h *history) Close() error {
	return nil
}

// AddPostWriteHook accepts a function to be run whenever a write to the buffer
// succeeds.
func (h *history) AddPostWriteHook(f func(string) error) {
	h.postWriteHooks = append(h.postWriteHooks, f)
}

// NewHistory returns a new history buffer.
func NewHistory(max int) *history {
	return &history{
		curr:  0,
		max:   max,
		lines: []string{},
	}
}
