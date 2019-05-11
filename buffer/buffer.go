// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package buffer

import (
	"strings"
	"time"

	"github.com/juju/loggo"
)

var log = loggo.GetLogger("stimmtausch.buffer")

// BufferLine represents a timestamped line of text.
type BufferLine struct {
	Timestamp time.Time
	Text      string
}

// Buffer represents a rolling buffer of lines used for input and output.
type Buffer struct {
	// The index of the current line
	curr int

	// The maximum number of lines to hold in the buffer.
	max int

	// The lines we're keeping track of.
	lines []*BufferLine

	// A list of functions to execute whenever the buffer is written to.
	postWriteHooks []func(*BufferLine) error
}

// add appends a line to the buffer, rolling a line out if necessary.
func (b *Buffer) add(line string) {
	l := &BufferLine{
		Timestamp: time.Now(),
		Text:      line,
	}
	b.lines = append(b.lines, l)
	if len(b.lines) > b.max {
		b.lines = b.lines[1 : b.max+1]
	}
	b.curr = len(b.lines) - 1
}

// Current returns the current line in the buffer.
func (b *Buffer) Current() *BufferLine {
	if len(b.lines) == 0 {
		return nil
	}
	return b.lines[b.curr]
}

// Forward moves the cursor forward in time one line and returns the current
// line's content.
func (b *Buffer) Forward() *BufferLine {
	b.curr++
	if b.curr >= len(b.lines) {
		b.curr = len(b.lines) - 1
	}
	return b.Current()
}

// Back moves the cursor back in time one line and returns the current
// line's content.
func (b *Buffer) Back() *BufferLine {
	if b.curr < 0 {
		b.curr = 0
	}
	line := b.Current()
	b.curr--
	return line
}

// Last moves the cursor to the last line.
func (b *Buffer) Last() *BufferLine {
	b.curr = len(b.lines) - 1
	return b.Current()
}

// OnLast returns true if the current line is the last (most recent) line.
func (b *Buffer) OnLast() bool {
	return b.curr == len(b.lines)-1
}

// String outputs the entire buffer as it stands.
func (b *Buffer) String() string {
	var builder strings.Builder
	for _, l := range b.lines {
		builder.WriteString(l.Text)
	}
	return builder.String()
}

// Write accepts a byte array and appends it to the buffer. It then executes
// every post-write hook. It returns the number of bytes written and any
// errors that occured.
// Fulfills io.Writer
func (b *Buffer) Write(line []byte) (int, error) {
	log.Tracef("received %d bytes", len(line))
	b.add(string(line))

	// This currently happens synchronously and serially. Should make sure we
	// want to do it this way or use a context.
	for _, hook := range b.postWriteHooks {
		err := hook(b.Current())
		if err != nil {
			return len(line), err
		}
	}
	return len(line), nil
}

// Size returns the number of lines in the buffer.
func (b *Buffer) Size() int {
	return len(b.lines)
}

// Close does nothing, and does it splendidly.
// Fulfills io.Closer
func (b *Buffer) Close() error {
	return nil
}

// AddPostWriteHook accepts a function to be run whenever a write to the buffer
// succeeds.
func (b *Buffer) AddPostWriteHook(f func(*BufferLine) error) {
	b.postWriteHooks = append(b.postWriteHooks, f)
}

// New returns a new buffer.
func New(max int) *Buffer {
	return &Buffer{
		curr:  0,
		max:   max,
		lines: []*BufferLine{},
	}
}
