// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

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

	// A list of triggers to match on input.
	Triggers []Trigger

	// References to compiled triggers.
	CompiledTriggers []*Trigger `yaml:"-" toml:"-"`

	Client Client

	HomeDir    string `yaml:"-" toml:"-"`
	ConfigDir  string `yaml:"-" toml:"-"`
	WorkingDir string `yaml:"-" toml:"-"`
	LogDir     string `yaml:"-" toml:"-"`
}

func (c *Config) FinalizeAndValidate() []error {
	log.Debugf("finalizing and validating config")
	var errs []error

	if c.Version == 0 {
		errs = append(errs, fmt.Errorf("version key wasn't set, perhaps no global configuration was loaded?"))
	}

	log.Tracef("finalizing and validating worlds")
	for name, world := range c.Worlds {
		world.Name = name
		if _, ok := c.Servers[world.Server]; !ok {
			errs = append(errs, fmt.Errorf("world %s refers to unknown server %s", name, world.Server))
		}
		c.Worlds[name] = world
	}

	log.Tracef("finalizing and validating servers")
	for name, server := range c.Servers {
		server.Name = name
		if _, ok := c.ServerTypes[server.ServerType]; server.ServerType != "" && !ok {
			errs = append(errs, fmt.Errorf("server %s refers to unknown server type %s", name, server.ServerType))
		}
		c.Servers[name] = server
	}

	log.Tracef("finalizing and validating triggers")
	for _, trigger := range c.Triggers {
		triggerRef, err := compileTrigger(trigger)
		if err != nil {
			errs = append(errs, err)
		}
		c.CompiledTriggers = append(c.CompiledTriggers, triggerRef)
	}

	c.HomeDir = HomeDir
	c.ConfigDir = ConfigDir
	c.WorkingDir = WorkingDir
	c.LogDir = LogDir
	return errs
}

// load populates a config object with configuration data from all available
// sources.
func load() (*Config, error) {
	var wrap wrapper
	snoot := snuffler.New(&wrap)

	log.Tracef("loading global config dirs")
	for _, location := range globalConfig {
		snoot.AddGlob(location)
	}

	log.Tracef("loading local config dirs")
	snoot.AddGlob(filepath.Join(ConfigDir, "*.st.*"))
	snoot.AddGlob(filepath.Join(ConfigDir, "*", "*.st.*"))

	if err := snoot.Snuffle(); err != nil {
		return nil, err
	}

	cfg := wrap.Stimmtausch
	errs := cfg.FinalizeAndValidate()
	if len(errs) != 0 {
		return nil, fmt.Errorf("could not validate config file: %+q", errs)
	}

	return &cfg, nil
}

// Reload loads config into a new object and then copies over its internals into
// the old object.
func (c *Config) Reload() error {
	log.Debugf("reloading configuration")

	newCfg, err := load()
	if err != nil {
		return err
	}

	c.Version = newCfg.Version
	c.ServerTypes = newCfg.ServerTypes
	c.Servers = newCfg.Servers
	c.Worlds = newCfg.Worlds
	c.Triggers = newCfg.Triggers
	c.CompiledTriggers = newCfg.CompiledTriggers
	c.Client = newCfg.Client
	return nil
}

// New returns a new config object all loaded and validated.
func New() (*Config, error) {
	log.Debugf("loading configuration")
	InitDirs()

	return load()
}
