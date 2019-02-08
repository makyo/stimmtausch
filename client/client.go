// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package client

import (
	"fmt"

	"github.com/makyo/stimmtausch/config"
)

// Client contains all of the information and objects Stimmtausch knows about.
// This pretty efficiently maps to information in the config file, and it may
// be worth simplifying that in the future.
type Client struct {
	// The current configuration object holding all worlds, servers, etc.
	config *config.Config

	// All active connections.
	connections map[string]*connection
}

// connectToWorld takes a given world and a connection name and creates a new
// connection in the client by calling connect on that world.
func (c *Client) connectToWorld(connectStr string, w config.World) (*connection, error) {
	log.Debugf("connecting to world %s (%s)", w.Name, connectStr)
	conn, err := NewConnection(connectStr, w, c.config.Servers[w.Server], c.config)
	if err != nil {
		log.Errorf("error connecting to world %s. %v", w.Name, err)
		return nil, err
	}
	c.connections[connectStr] = conn

	return conn, nil
}

// connectToServer will connect to a server with a new world created on the spot
// for that purpose.
func (c *Client) connectToServer(connectStr string, s config.Server) (*connection, error) {
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
	w, ok := c.config.Worlds[connectStr]
	if ok {
		conn, err := c.connectToWorld(connectStr, w)
		if err != nil {
			log.Errorf("unable to connect to world %s: %v", connectStr, err)
			return nil, err
		}
		return conn, nil
	}

	log.Tracef("checking if it's a server...")
	s, ok := c.config.Servers[connectStr]
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

// New creates a new Client and populates it using information from the config.
func New(cfg *config.Config) (*Client, error) {
	log.Debugf("creating client")
	c := &Client{
		config:      cfg,
		connections: map[string]*connection{},
	}
	return c, nil
}
