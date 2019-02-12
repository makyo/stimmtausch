// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.
//
// +build !windows

package main

import (
	"os"
	"path/filepath"

	"github.com/juju/loggo"

	"github.com/makyo/stimmtausch/cmd"
	"github.com/makyo/stimmtausch/config"
)

func main() {
	config.InitDirs()
	f, err := os.Create(filepath.Join(config.LogDir, "stimmtausch.log"))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	loggo.ReplaceDefaultWriter(loggo.NewSimpleWriter(f, loggo.DefaultFormatter))

	cmd.Execute()
}
