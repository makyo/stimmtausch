// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package client

import (
	"path/filepath"

	"github.com/makyo/st/config"
)

// world represents the union between a server and a character.
type world struct {
	// The key for the world in the configuration file.
	name string

	// A free-form display name for the world.
	displayName string

	// The server to which this world belongs.
	server *server

	// The username and password to connect with.
	username, password string

	// Whether or not to maintain a rotated log of each connection to this world.
	log bool
}

// getWorldFile returns a file (or directory) name within the scope of the
// world. These live in $HOME/.config/stimmtausch/worlds/{worldname}.
func (w *world) getWorldFile(name string) (string, error) {
	home, err := config.HomeDir()
	if err != nil {
		panic(err)
	}
	filename := filepath.Join(home, "worlds", w.name, name)
	log.Tracef("file path: %s", filename)
	return filename, nil
}

// connect creates a new connection for the world with the provided name.
func (w *world) connect(name string) (*connection, error) {
	return NewConnection(name, w)
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
