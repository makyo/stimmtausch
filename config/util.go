// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package config

import (
	"fmt"
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
)

// InitDirs initializes the directories used by Stimmtausch.
func InitDirs() {
	if HomeDir != "" && ConfigDir != "" && WorkingDir != "" && LogDir != "" {
		return
	}
	HomeDir, err := homedir.Dir()
	if err != nil {
		panic(fmt.Sprintf("could not find home dir: %v", err))
	}
	startDir := os.Getenv("SNAP_USER_COMMON")
	if startDir == "" {
		WorkingDir = filepath.Join(HomeDir, ".local", "share", "stimmtausch")
		LogDir = filepath.Join(HomeDir, ".local", "log", "stimmtausch")
	} else {
		WorkingDir = filepath.Join(startDir, "worlds")
		LogDir = filepath.Join(startDir, "logs")
	}
	ConfigDir = filepath.Join(HomeDir, ".config", "stimmtausch")
	globalConfig = []string{
		// Locations for installed configuration files.
		"/etc/stimmtausch/st.yaml",
		"/etc/stimmtausch/conf.d/*.yaml",
		"/etc/stimmtausch/conf.d/*.toml",
		"/etc/stimmtausch/conf.d/*.json",

		// Locations for development configuration files.
		"_conf/global/st.yaml",
		"_conf/global/conf.d/*.yaml",
		"_conf/global/conf.d/*.toml",
		"_conf/global/conf.d/*.json",
		"_conf/local/*.yaml",
		"_conf/local/*.toml",
		"_conf/local/*.json",
	}
}
