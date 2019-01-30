package client

import (
	"regexp"
)

type trigger interface {
	Run(string) string
}

type hilite struct {
	filter *regexp.Regexp
}

type gag struct {
	filter *regexp.Regexp
}

type script struct {
	filter *regexp.Regexp
}
