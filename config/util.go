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
	globalConfigDirs   []string
	globalMasterConfig string
	HomeDir            string
	ConfigDir          string
	WorkingDir         string
	LogDir             string
	Environment        string
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
	if Environment := os.Getenv("ST_ENV"); !(Environment == "DEV" || Environment == "PROD") {
		Environment = "PROD"
	}
	if Environment == "PROD" {
		globalConfigDirs = []string{"/etc/stimmtausch/conf.d/*"}
		globalMasterConfig = "/etc/stimmtausch/st.yaml"
	} else {
		globalConfigDirs = []string{"_conf/conf.d/*"}
		globalMasterConfig = "_conf/st.yaml"
	}
	return nil
}
