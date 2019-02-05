package config

import (
	"github.com/makyo/snuffler"
)

var log = loggo.GetLogger("stimmtausch.config")

// wrapper just wraps a Config object, since, for readability's sake, all
// Stimmtausch yaml files have everything under the `stimmtausch` key.
type wrapper struct {
	Stimmtausch Config
}

// Config holds all the configuration for Stimmtausch.
type Config struct {
	Version     int                   `yaml:",omitempty"`
	Client      Client                `yaml:",omitempty"`
	ServerTypes map[string]ServerType `yaml:"server_types,omitempty" toml:"server_types"`
	Servers     map[string]Server     `yaml:",omitempty"`
	Worlds      map[string]World      `yaml:",omitempty"`
}

// Client holds information about various ways in which the client/ui act.
type Client struct {
	Syslog  Syslog  `yaml:",omitempty"`
	Logging Logging `yaml:",omitempty"`
}

// Syslog holds information about system logging (rather than world logging).
type Syslog struct {
	ShowSyslog bool   `yaml:"show_syslog" toml:"log_level"`
	LogLevel   string `yaml:"log_level" toml:"log_level"`
}

// Logging holds information about world logging (rather than system logging).
type Logging struct {
	TimeString    string `yaml:"time_string" toml:"time_string"`
	LogTimestamps bool   `yaml:"log_timestamps" toml"log_timestamps"`
	LogWorld      bool   `yaml:"log_world" toml:"log_world"`
}

func Load(additionalLocations []string) (*Config, error) {
	var cfg wrapper
	snoot := snuffler.New(&cfg)

	if err := snoot.AddFile(globalMasterConfig); err != nil {
		return nil, err
	}
	for _, location := range globalConfigDirs {
		if err := snoot.AddGlob(location); err != nil {
			return nil, err
		}
	}
	for _, location := range additionalLocations {
		if err := snoot.AddGlob(location); err != nil {
			return nil, err
		}
	}

	if err := snoot.Snuffle(); err != nil {
		return nil, err
	}

	return cfg.Stimmtausch
}
