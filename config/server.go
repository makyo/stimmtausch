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
	ServerType string `yaml:"type" toml:"type"`
}

// ServerType represents a type of server (MUCK, MUSH, etc...), which mostly
// boils down to things such as how to connect to it, etc.
type ServerType struct {
	Name             string
	ConnectString    string `yaml:"connect_string" toml:"connect_string"`
	DisconnectString string `yaml:"disconnect_string" toml:"disconnect_string"`
}

// NewServer returns a new server object for the given values.
func NewServer(name, host string, port uint, ssl, insecure bool, srvType string) *Server {
	return &Server{
		Name:       name,
		Host:       host,
		Port:       port,
		SSL:        ssl,
		Insecure:   insecure,
		ServerType: srvType,
	}
}

// NewServerType returns a new serverType object for the given values.
func NewServerType(name, connectString, disconnectString string) *ServerType {
	return &ServerType{
		Name:             name,
		ConnectString:    connectString,
		DisconnectString: disconnectString,
	}
}
