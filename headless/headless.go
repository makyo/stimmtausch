// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package headless

import (
	"os"

	"github.com/juju/loggo"

	"github.com/makyo/stimmtausch/config"
)

var log = loggo.GetLogger("stimmtausch.headless")

func New(args []string, cfg *config.Config) {
	log.Criticalf("not implemented")
	os.Exit(2)
}
