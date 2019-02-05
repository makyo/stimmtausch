// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.
//
// +build !windows

package config

import (
	"path/filepath"

	"github.com/juju/loggo"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	globalConfigDirs   = []string{"/etc/stimmtausch/conf.d/*"}
	globalMasterConfig = "/etc/stimmtausch/st.yaml"
)

// ConfigDir returns the directory in which Stimmtausch expects its config
// files to be stored.
func configDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		log.Criticalf("could not find homedir: %v", err)
		return "", err
	}
	return filepath.Join(home, ".config", "stimmtausch"), nil
}

// WorkingDir returns the directory in which Stimmtausch does all of its work.
func WorkingDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		log.Criticalf("could not find homedir: %v", err)
		return "", err
	}
	return filepath.Join(home, ".local", "share", "stimmtausch"), nil
}
