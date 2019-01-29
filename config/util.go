package config

import (
	"path/filepath"

	"github.com/juju/loggo"
	homedir "github.com/mitchellh/go-homedir"
)

var log = loggo.GetLogger("stimmtausch.config")

func HomeDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		log.Criticalf("could not find homedir: %v", err)
		return "", err
	}
	return filepath.Join(home, ".config", "stimmtausch"), nil
}
