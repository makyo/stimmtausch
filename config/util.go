// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package config

import (
	"path/filepath"

	"github.com/juju/loggo"
	homedir "github.com/mitchellh/go-homedir"
)

var log = loggo.GetLogger("stimmtausch.config")

// HomeDir returns the directory in which Stimmtausch does all of its work.
func HomeDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		log.Criticalf("could not find homedir: %v", err)
		return "", err
	}
	return filepath.Join(home, ".config", "stimmtausch"), nil
}
