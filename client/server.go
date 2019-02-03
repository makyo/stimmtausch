// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package client

type server struct {
	name       string
	host       string
	port       uint
	ssl        bool
	insecure   bool
	serverType *serverType
}

type serverType struct {
	name          string
	connectString string
}

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

func NewServerType(name, connectString string) *serverType {
	return &serverType{
		name:          name,
		connectString: connectString,
	}
}
