// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package config

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
