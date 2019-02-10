// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package config

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

var (
	globalConfig []string
	HomeDir      string
	ConfigDir    string
	WorkingDir   string
	LogDir       string
	Environment  string
)

// initDirs initializes the directories used by Stimmtausch.
func initEnv() error {
	HomeDir, err := homedir.Dir()
	if err != nil {
		log.Criticalf("could not find homedir: %v", err)
		return err
	}
	ConfigDir = filepath.Join(HomeDir, ".config", "stimmtausch")
	WorkingDir = filepath.Join(HomeDir, ".local", "share", "stimmtausch")
	LogDir = filepath.Join(HomeDir, ".local", "log", "stimmtausch")
	globalConfig = []string{
		// Locations for installed configuration files.
		"/etc/stimmtausch/st.yaml",
		"/etc/stimmtausch/conf.d/*.yaml",
		"/etc/stimmtausch/conf.d/*.toml",
		"/etc/stimmtausch/conf.d/*.json",

		// Locations for development configuration files.
		"_conf/st.yaml",
		"_conf/conf.d/*.yaml",
		"_conf/conf.d/*.toml",
		"_conf/conf.d/*.json",
	}
	return nil
}
