// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package config

// World represents the union between a server and a character.
type World struct {
	// The key for the world in the configuration file.
	Name string

	// A free-form display name for the world.
	DisplayName string `yaml:"display_name" toml:"display_name"`

	// The server to which this world belongs.
	Server string

	// The username and password to connect with.
	Username, Password string

	// Whether or not to maintain a rotated log of each connection to this world.
	Log bool
}

// NewWorld returns a new world object for the given values.
func NewWorld(name, displayName, srv, username, password string, logByDefault bool) *World {
	return &World{
		Name:        name,
		DisplayName: displayName,
		Server:      srv,
		Username:    username,
		Password:    password,
		Log:         logByDefault,
	}
}
