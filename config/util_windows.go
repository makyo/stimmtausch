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

var (
	globalConfigDirs   = []string{filepath.Join(os.Getenv("ProgramFiles"), "Stimmtausch", "config", "*")}
	globalMasterConfig = filepath.Join(os.Getenv("ProgramFiles"), "Stimmtausch", "st.yaml")
)

// ConfigDir returns the directory in which Stimmtausch expects its config
// files to be stored.
func configDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		log.Criticalf("could not find homedir: %v", err)
		return "", err
	}
	return filepath.Join(home, "AppData", "Stimmtausch", "Configuration"), nil
}

// WorkingDir returns the directory in which Stimmtausch does all of its work.
func WorkingDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		log.Criticalf("could not find homedir: %v", err)
		return "", err
	}
	return filepath.Join(home, "AppData", "Stimmtausch", "Working Directory"), nil
}
