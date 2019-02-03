// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package client

import (
	"fmt"

	"github.com/spf13/viper"
)

// Client contains all of the information and objects Stimmtausch knows about.
// This pretty efficiently maps to information in the config file, and it may
// be worth simplifying that in the future.
type Client struct {
	// All active connections.
	connections map[string]*connection

	// All worlds the client knows about, connected or not.
	worlds map[string]*world

	// All servers, whether or not worlds are using them.
	servers map[string]*server

	// All server types, whether or not they're used.
	serverTypes map[string]*serverType

	// The default server type.
	defaultServerType string

	// The default world to connect to.
	defaultWorld string
}

// New creates a new Client and populates it using information from the config.
func New() (*Client, error) {
	log.Debugf("creating client")
	c := &Client{
		worlds:            map[string]*world{},
		servers:           map[string]*server{},
		serverTypes:       map[string]*serverType{},
		connections:       map[string]*connection{},
		defaultServerType: viper.GetString("stimmtausch.default_server_type"),
		defaultWorld:      viper.GetString("stimmtausch.default_world"),
	}

	log.Tracef("loading server types")
	for serverTypeName, spec := range viper.GetStringMap("stimmtausch.server_types") {
		s, ok := spec.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed server_types entry in config")
		}
		if err := c.UpsertServerType(serverTypeName, s); err != nil {
			return nil, err
		}
	}

	log.Tracef("loading servers")
	for serverName, spec := range viper.GetStringMap("stimmtausch.servers") {
		s, ok := spec.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed servers entry in config")
		}
		if err := c.UpsertServer(serverName, s); err != nil {
			return nil, err
		}
	}

	log.Tracef("loading worlds")
	for worldName, spec := range viper.GetStringMap("stimmtausch.worlds") {
		s, ok := spec.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed worlds entry in config")
		}
		if err := c.UpsertWorld(worldName, s); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// UpsertServerType inserts a new or updates an existing server type.
func (c *Client) UpsertServerType(name string, spec map[string]interface{}) error {
	log.Debugf("upserting server type %s", name)
	var connectString string

	if inter, ok := spec["connect_string"]; ok {
		connectString, ok = inter.(string)
		if !ok {
			return fmt.Errorf("malformed connect_string for server type %s", name)
		}
	} else {
		return fmt.Errorf("server type %s missing connect_string key", name)
	}

	if _, ok := c.serverTypes[name]; ok {
		log.Infof("updating server type %s", name)
	} else {
		log.Infof("creating server type %s", name)
	}

	c.serverTypes[name] = NewServerType(name, connectString)

	return nil
}

// UpsertServer inserts a new or updates an existing server.
func (c *Client) UpsertServer(name string, spec map[string]interface{}) error {
	log.Debugf("upserting server %s", name)
	var host string
	var port int
	ssl := false
	insecure := false
	var st *serverType

	if inter, ok := spec["host"]; ok {
		host, ok = inter.(string)
		if !ok {
			return fmt.Errorf("malformed host key for server %s", name)
		}
	} else {
		return fmt.Errorf("server %s missing host key", name)
	}

	if inter, ok := spec["port"]; ok {
		port, ok = inter.(int)
		if !ok {
			return fmt.Errorf("malformed port key for server %s", name)
		}
	} else {
		return fmt.Errorf("server %s missing port key", name)
	}

	if inter, ok := spec["ssl"]; ok {
		ssl, ok = inter.(bool)
		if !ok {
			return fmt.Errorf("malformed ssl key for server %s", name)
		}
	}

	if inter, ok := spec["insecure"]; ok {
		insecure, ok = inter.(bool)
		if !ok {
			return fmt.Errorf("malformed insecure key for server %s", name)
		}
	}

	if inter, ok := spec["type"]; ok {
		stStr, ok := inter.(string)
		if !ok {
			return fmt.Errorf("malformed type for server %s", name)
		}
		st, ok = c.serverTypes[stStr]
		if !ok {
			return fmt.Errorf("server type %s for server %s refers to a type that doesn't exist", stStr, name)
		}
	} else {
		st, ok = c.serverTypes[c.defaultServerType]
		if !ok {
			return fmt.Errorf("default_server_type %s refers to a type that doesn't exist", c.defaultServerType)
		}
	}

	if _, ok := c.servers[name]; ok {
		log.Infof("updating server %s", name)
	} else {
		log.Infof("creating server %s", name)
	}

	c.servers[name] = NewServer(name, host, uint(port), ssl, insecure, st)

	return nil
}

// UpsertWorld inserts a new or updates an existing world.
func (c *Client) UpsertWorld(name string, spec map[string]interface{}) error {
	log.Debugf("upserting world %s", name)
	var srv *server
	var displayName string
	var username string
	var password string
	logByDefault := false

	if inter, ok := spec["server"]; ok {
		srvName, ok := inter.(string)
		if !ok {
			return fmt.Errorf("malformed server for world %s", name)
		}
		if srv, ok = c.servers[srvName]; !ok {
			return fmt.Errorf("world %s references undefined server %s", name, srvName)
		}
	} else {
		return fmt.Errorf("world %s missing server key", name)
	}

	if inter, ok := spec["display_name"]; ok {
		displayName, ok = inter.(string)
		if !ok {
			return fmt.Errorf("malformed display_name for world %s", name)
		}
	} else {
		displayName = name
	}

	if inter, ok := spec["username"]; ok {
		username, ok = inter.(string)
		if !ok {
			return fmt.Errorf("malformed username for world %s", name)
		}
	} else {
		return fmt.Errorf("world %s missing username key", name)
	}

	if inter, ok := spec["password"]; ok {
		password, ok = inter.(string)
		if !ok {
			return fmt.Errorf("malformed password for world %s", name)
		}
	} else {
		return fmt.Errorf("world %s missing password key", name)
	}

	if inter, ok := spec["log"]; ok {
		logByDefault, ok = inter.(bool)
		if !ok {
			return fmt.Errorf("malformed log key for world %s", name)
		}
	}

	if _, ok := c.worlds[name]; ok {
		log.Infof("updating world %s", name)
	} else {
		log.Infof("creating world %s", name)
	}

	c.worlds[name] = NewWorld(name, displayName, srv, username, password, logByDefault)

	return nil
}

// connectToWorld takes a given world and a connection name and creates a new
// connection in the client by calling connect on that world.
func (c *Client) connectToWorld(connectStr string, w *world) (*connection, error) {
	log.Debugf("connecting to world %s (%s)", w.name, connectStr)
	conn, err := w.connect(connectStr)
	if err != nil {
		log.Errorf("error connecting to world %s. %v", w.name, err)
		return nil, err
	}
	c.connections[connectStr] = conn

	return conn, nil
}

// connectToServer will connect to a server with a new world created on the spot
// for that purpose.
func (c *Client) connectToServer(connectStr string, s *server) (*connection, error) {
	return nil, fmt.Errorf("not implemented")
}

// connectToRaw will attempt to connect to a host:port string, building a
// server and world for the purpose.
func (c *Client) connectToRaw(connectStr string) (*connection, error) {
	return nil, fmt.Errorf("not implemented")
}

// Connect accepts a string and tries to connect to it in the following ways:
// * If the string is the name of a world, it will connect that world; otherwise
// * If the string is the name of a server, it will connect to it with a new
//   world created for that purpose; otherwise
// * It will try to connect to that string as if it were a host:port; finally
// * It will fail.
func (c *Client) Connect(connectStr string) (*connection, error) {
	log.Debugf("attempting to connect to %s in %v", connectStr, c)

	log.Tracef("checking if it's a world...")
	w, ok := c.worlds[connectStr]
	if ok {
		conn, err := c.connectToWorld(connectStr, w)
		if err != nil {
			log.Errorf("unable to connect to world %s: %v", connectStr, err)
			return nil, err
		}
		return conn, nil
	}

	log.Tracef("checking if it's a server...")
	s, ok := c.servers[connectStr]
	if ok {
		conn, err := c.connectToServer(connectStr, s)
		if err != nil {
			log.Errorf("unable to connect to server %s: %v", connectStr, err)
			return nil, err
		}
		return conn, nil
	}

	log.Tracef("defaulting to trying it as an address")
	conn, err := c.connectToRaw(connectStr)
	if err != nil {
		log.Errorf("unable to connect to address %s: %v", connectStr, err)
		return nil, err
	}
	return conn, nil
}

// Close will close a connection with the given name (usually the connectStr).
func (c *Client) Close(name string) {
	log.Debugf("closing connection %s", name)
	c.connections[name].Close()
}

// CloseAll will attempt to close all open connections.
func (c *Client) CloseAll() {
	log.Debugf("closing all connections")
	for _, conn := range c.connections {
		conn.Close()
	}
}
