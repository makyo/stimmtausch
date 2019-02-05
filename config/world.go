// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package config

import (
	"path/filepath"
)

// World represents the union between a server and a character.
type World struct {
	// The key for the world in the configuration file.
	Name string

	// A free-form display name for the world.
	DisplayName string

	// The server to which this world belongs.
	ServerName string

	// The username and password to connect with.
	Username, Password string

	// Whether or not to maintain a rotated log of each connection to this world.
	Log bool
}

// GetWorldFile returns a file (or directory) name within the scope of the
// world. These live in $HOME/.config/stimmtausch/worlds/{worldname}.
func (w *World) GetWorldFile(name string) (string, error) {
	home, err := config.HomeDir()
	if err != nil {
		panic(err)
	}
	filename := filepath.Join(home, "worlds", w.name, name)
	log.Tracef("file path: %s", filename)
	return filename, nil
}

// NewWorld returns a new world object for the given values.
func NewWorld(name, displayName string, srv *server, username, password string, logByDefault bool) *world {
	return &world{
		name:        name,
		displayName: displayName,
		server:      srv,
		username:    username,
		password:    password,
		log:         logByDefault,
	}
}
