// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package help

import (
	"github.com/juju/loggo"
)

var log = loggo.GetLogger("stimmtausch.help")

type Help struct {
	Name        string
	ShortDesc   string
	Synopsis    map[string]string
	Overview    string
	Description string
	SeeAlso     string
}
