package config

import (
	"fmt"
	"path/filepath"

	"github.com/juju/loggo"

	"github.com/makyo/snuffler"
)

var log = loggo.GetLogger("stimmtausch.config")

// wrapper just wraps a Config object, since, for readability's sake, all
// Stimmtausch yaml files have everything under the `stimmtausch` key.
type wrapper struct {
	Stimmtausch Config
}

// Config holds all the configuration for Stimmtausch.
// For information about what the settings are used for and how they should
// appear, see the files in _conf
type Config struct {
	// Version of the configuration structure.
	Version int

	// A list of server types (MUCK, MUD, etc), dictating how to connect.
	ServerTypes map[string]ServerType `yaml:"server_types" toml:"server_types"`

	// A list of servers.
	Servers map[string]Server

	// A list of worlds (characters/users/accounts/etc) tying login information
	// to servers.
	Worlds map[string]World

	// Information regarding how Stimmtausch runs.
	Client struct {

		// Information regarding the logging generated by the program (as
		// opposed to the connections).
		Syslog struct {

			// Whether or not to show the system log in a pane in the UI.
			ShowSyslog bool `yaml:"show_syslog" toml:"log_level"`

			// Lowest log level to show by default. Options are:
			// TRACE, DEBUG, INFO*, WARNING, ERROR, CRITICAL
			LogLevel string `yaml:"log_level" toml:"log_level"`
		}

		// Information regarding logging from connections.
		Logging struct {

			// The time format to use when logging times.
			TimeString string `yaml:"time_string" toml:"time_string"`

			// Whether or not to log timestamps.
			LogTimestamps bool `yaml:"log_timestamps" toml"log_timestamps"`

			// Whether or not to keep logs of the connection after disconnect.
			LogWorld bool `yaml:"log_world" toml:"log_world"`
		}

		// Information regarding the user interface.
		UI struct {

			// How many lines of scrollback (data received) to keep in memory.
			Scrollback int

			// How many lines of history (data sent) to keep in memory.
			History int

			// Whether or not to use a unified history buffer for all
			// connections, or one per.
			UnifiedHistoryBuffer bool
		}
	}
}

func (c *Config) finalizeAndValidate() []error {
	log.Debugf("finalizing and validating config")
	var errs []error

	log.Tracef("finalizing and validating worlds")
	for name, world := range c.Worlds {
		world.Name = name
		if _, ok := c.Servers[world.Server]; !ok {
			errs = append(errs, fmt.Errorf("world %s refers to unknown server %s", name, world.Server))
		}
	}

	log.Tracef("finalizing and validating servers")
	for name, server := range c.Servers {
		server.Name = name
		if _, ok := c.ServerTypes[server.ServerType]; server.ServerType != "" && !ok {
			errs = append(errs, fmt.Errorf("server %s refers to unknown server type %s", name, server.ServerType))
		}
	}
	return errs
}

// Load populates a config object with configuration data from all available
// sources.
func Load(additionalLocations []string) (*Config, error) {
	log.Debugf("loading configuration")
	err := initEnv()
	if err != nil {
		return nil, err
	}

	var wrap wrapper
	snoot := snuffler.New(&wrap)

	log.Tracef("loading global master config")
	if err := snoot.AddFile(globalMasterConfig); err != nil {
		return nil, err
	}

	log.Tracef("loading global config dirs")
	for _, location := range globalConfigDirs {
		snoot.AddGlob(location)
	}

	log.Tracef("loading local config dirs")
	snoot.AddGlob(filepath.Join(ConfigDir, "*.st.*"))
	snoot.AddGlob(filepath.Join(ConfigDir, "*", "*.st.*"))
	snoot.MaybeAddFile(filepath.Join(HomeDir, ".strc"))

	log.Tracef("loading additional locations")
	for _, location := range additionalLocations {
		snoot.AddGlob(location)
	}

	if err := snoot.Snuffle(); err != nil {
		return nil, err
	}

	cfg := wrap.Stimmtausch
	errs := cfg.finalizeAndValidate()
	if len(errs) != 0 {
		return nil, fmt.Errorf("could not validate config file: %+q", errs)
	}

	return &cfg, nil
}