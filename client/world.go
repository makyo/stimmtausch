package client

import (
	"path/filepath"

	"github.com/makyo/st/config"
)

type world struct {
	name        string
	displayName string
	server      *server
	username    string
	password    string
	log         bool
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

func (w *world) connect(name string) (*connection, error) {
	return NewConnection(name, w)
}

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
