package ui

import (
	"strings"
)

type history struct {
	curr           int
	max            int
	lines          []string
	postWriteHooks []func() error
}

func (h *history) add(line string) {
	h.lines = append(h.lines, strings.TrimSpace(line))
	if len(h.lines) > h.max {
		h.lines = h.lines[:h.max]
	}
	h.curr = len(h.lines) - 1
}

func (h *history) current() string {
	return h.lines[h.curr]
}

func (h *history) forward() string {
	h.curr++
	if h.curr >= len(h.lines) {
		h.curr = len(h.lines) - 1
	}
	return h.current()
}

func (h *history) back() string {
	if h.curr < 0 {
		h.curr = 0
	}
	line := h.current()
	h.curr--
	return line
}

func (h *history) onLast() bool {
	return h.curr == len(h.lines)-1
}

func (h *history) String() string {
	return strings.TrimSpace(strings.Join(h.lines, ""))
}

func (h *history) Write(line []byte) (int, error) {
	h.add(string(line))
	return len(line), nil
}

func (h *history) Close() error {
	return nil
}

func (h *history) AddPostWriteHook(f func() error) {
	h.postWriteHooks = append(h.postWriteHooks, f)
}

func NewHistory(max int) *history {
	return &history{
		curr:  0,
		max:   max,
		lines: []string{},
	}
}
