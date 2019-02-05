// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package config

// Server represents information required to connect to a remote server.
type Server struct {
	// The key for the server in the configuration file.
	Name string

	// The hostname for the server.
	Host string

	// The port to connect to.
	Port uint

	// Whether or not to use SSL.
	SSL bool

	// Whether or not self-signed certs should be trusted.
	Insecure bool

	// The type of server (MUCK, MUSH, etc...) this is.
	ServerType string
}

// ServerType represents a type of server (MUCK, MUSH, etc...), which mostly
// boils down to things such as how to connect to it, etc.
type ServerType struct {
	Name             string
	ConnectString    string
	DisconnectString string
}

// NewServer returns a new server object for the given values.
func NewServer(name, host string, port uint, ssl, insecure bool, srvType *serverType) *server {
	return &server{
		name:       name,
		host:       host,
		port:       port,
		ssl:        ssl,
		insecure:   insecure,
		serverType: srvType,
	}
}

// NewServerType returns a new serverType object for the given values.
func NewServerType(name, connectString string) *serverType {
	return &serverType{
		name:          name,
		connectString: connectString,
	}
}
