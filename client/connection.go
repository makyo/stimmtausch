// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/juju/loggo"

	"github.com/makyo/st/config"
)

var log = loggo.GetLogger("stimmtausch.client")

// Hardcoded client settings.
const (
	// The name of the FIFO file.
	inFile string = "in"

	// The name of the global output file.
	outFile string = "out"

	// The size of buffer to read from the connection.
	bufferSize int = 1024

	// The delay used when reading from the FIFO. See note below.
	fifoReadDelay = 100 * time.Millisecond

	// The keepalive setting for connections. See note below.
	keepalive = 15 * time.Minute
)

// getTimestamp gets the current time in the format specified above.
func (c *connection) getTimestamp() string {
	return time.Now().Format(c.config.Client.Logging.TimeString)
}

// output represents a named io.WriteCloser.
type output struct {
	// A name used for logging and referencing down the line.
	name string

	// Whether or not this is the global output (XXX is this necessary?).
	global bool

	// The io.Writecloser itself.
	output io.WriteCloser
}

// conn stores all connection settings
type connection struct {
	// The name (usually connectStr) of the connection.
	name string

	// The world and server to which this connection belongs.
	// These are maintained separately from the app config as they may be
	// connected and passed in for settings not in the user's config.
	world  config.World
	server config.Server

	// The app configuration.
	config *config.Config

	// The TCP address of the server.
	addr *net.TCPAddr

	// The TCP connection itself.
	connection net.Conn

	// The FIFO file used for maintaining the connection.
	fifo *os.File

	// The array of io.WriteClosers that output from the world is written to.
	outputs []*output

	// A channel signalling a disconnect request.
	disconnect chan bool

	// A channel signalling that the server has disconnected.
	disconnected chan bool

	// Whether or not the server is connected.
	connected bool
}

// lookupHostname gets the TCP address for the world's hostname.
func (c *connection) lookupHostname() error {
	log.Debugf("attempting to resolve %s:%d", c.server.Host, c.server.Port)
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.server.Host, c.server.Port))
	if err != nil {
		log.Errorf("unable to resolve host %s: %v", c.server.Host, err)
		return err
	}
	c.addr = addr
	log.Debugf("host %s resolves to %v", c.server.Host, c.addr)
	return nil
}

// getConnectionFile returns a file (or directory) name within the scope of
// the connection. These live in
// $HOME/.local/share/stimmtausch/{connname}.
func (c *connection) getConnectionFile(name string) string {
	return filepath.Join(c.config.WorkingDir, c.name, name)
}

// getLogFile returns a file (or directory) name within the scope of the
// connection for the sake of logging. These live in
// $HOME/.local/log/stimmtausch/{worldname}.
func (c *connection) getLogFile(name string) string {
	return filepath.Join(c.config.LogDir, c.world.Name, name)
}

// makeFIFO creates the FIFO file for the world, used to manage the information
// sent to and recieved from the connection.
func (c *connection) makeFIFO() error {
	log.Debugf("creating FIFO file for %s", c.name)
	file := c.getConnectionFile(inFile)
	var err error

	log.Tracef("checking if FIFO exists")
	if _, err = os.Stat(file); err == nil {
		log.Criticalf("FIFO for connection %s already exists!", c.world.Name)
	}

	log.Tracef("making FIFO")
	if err = syscall.Mkfifo(file, 0644); err != nil {
		log.Criticalf("unable to make FIFO for %s!", c.world.Name, err)
		return err
	}
	log.Tracef("FIFO created as %s", file)

	log.Tracef("opening FIFO")
	if c.fifo, err = os.OpenFile(file, os.O_RDONLY|syscall.O_NONBLOCK, os.ModeNamedPipe); err != nil {
		log.Criticalf("unable to open FIFO for reading %s! %v", file, err)
		return err
	}
	log.Tracef("FIFO opened as %s", c.fifo.Name())
	return nil
}

// makeLogfile creates a logfile from a given name.
func (c *connection) makeLogfile(out *output) error {
	log.Debugf("creating a log file for %s", c.name)

	log.Tracef("checking if %s exists", out.name)
	_, err := os.Stat(out.name)
	if err == nil {
		fmt.Printf("Warning: %v already exists; appending.\n", out.name)
	}

	log.Tracef("opening %s for logging", c.name)
	f, err := os.OpenFile(out.name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Warningf("could not open logfile %s for %s, not logging. %v", out.name, c.world.Name, err)
		return err
	}

	out.output = f
	log.Debugf("logfile created as %s", out.name)
	return nil
}

// connect creates a TCP connection to the world's TCP address. It connects
// over SSL if available and, if requested, will use insecure certs.
func (c *connection) connect() error {
	log.Debugf("creating TCP connection for %s", c.name)
	var err error
	if c.server.SSL {
		log.Tracef("creating SSL connection")
		var conf *tls.Config
		if c.server.Insecure {
			conf = &tls.Config{InsecureSkipVerify: true}
		} else {
			conf = &tls.Config{ServerName: c.server.Host}
		}
		if c.connection, err = tls.Dial("tcp", c.addr.String(), conf); err != nil {
			log.Criticalf("unable to dial %v over SSL for %s! %v", c.addr, c.world.Name, err)
			return err
		}
		log.Debugf("connected to server over SSL for %s", c.world.Name)
	} else {
		log.Tracef("creating regular TCP connection")
		if c.connection, err = net.DialTCP("tcp", nil, c.addr); err != nil {
			log.Criticalf("unable to dial %v for %s! %v", c.addr, c.world.Name, err)
			return err
		}
		log.Debugf("connected to server for %s", c.world.Name)
	}
	//
	//		// XXX This doesn't work with SSL connections, need to find an alternative...
	//		log.Tracef("attempting to set a keepalive for %s", c.name)
	//		if err = c.connection.SetKeepAlive(true); err != nil {
	//			log.Warningf("unable to set keep alive for %s - you may get booted. %v", c.world.Name, err)
	//		}
	//		if err = c.connection.SetKeepAlivePeriod(keepalive); err != nil {
	//			log.Warningf("unable to set keep alive period for %s - you may get booted. %v", c.world.Name, err)
	//		}
	c.connected = true
	return nil
}

// readToConn reads from the FIFO and sends to the connection.
func (c *connection) readToConn() {
	tmpError := fmt.Sprintf("read %v: resource temporarily unavailable", c.fifo.Name())
	for {
		select {
		case <-c.disconnect:
			log.Debugf("%s received disconnect; returning", c.world.Name)
			c.disconnected <- true
			return
		default:
			// This pause between reads from the FIFO is the difference between 0.2%
			// and 100% cpu usage when idle. Also without this you will get excessive
			// "read %v: resource temporarily unavailable" errors on some OSes.
			time.Sleep(fifoReadDelay)
			buf := make([]byte, bufferSize)
			bytesIn, err := c.fifo.Read(buf)
			if err != nil && err.Error() != "EOF" && err.Error() != tmpError {
				log.Criticalf("FIFO broke??¿? World %s. %v", c.world.Name, err)
			} else if bytesIn == 0 {
				continue
			}
			log.Tracef("%d bytes read from FIFO", bytesIn)
			bytesOut, err := c.connection.Write(buf[:bytesIn])
			if err != nil {
				log.Criticalf("FIFO broke??¿? World %s. %v", c.world.Name, err)
			}
			log.Tracef("%d bytes written to connection", bytesOut)
		}
	}
}

// readToFile reads from the connection and writes to outfiles.
func (c *connection) readToFile() {
	reader := bufio.NewReader(c.connection)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			if !c.connected {
				return
			}
			log.Warningf("server disconnected with %v", err)
			disconnectMsg := fmt.Sprintf("\n~Connection lost at %v\n", c.getTimestamp())
			for _, out := range c.outputs {
				if _, err := fmt.Fprintln(out.output, disconnectMsg); err != nil {
					log.Warningf("unable to write to output %s for %s. %v", out.name, c.world.Name, err)
				}
			}
			c.Close()
			return
		}
		log.Tracef("%d characters read from %s", len(line), c.world.Name)

		for _, out := range c.outputs {
			bytesOut, err := fmt.Fprintln(out.output, line)
			if err != nil {
				log.Warningf("unable to write to output %s for world %s. %v", out.name, c.world.Name, err)
			}
			log.Tracef("%d bytes written to output %s for %s", bytesOut, out.name, c.world.Name)
		}
	}
}

// closeConnection closes the world's TCP connection.
func (c *connection) closeConnection() {
	if !c.connected {
		log.Debugf("%s already closed", c.name)
		return
	}
	log.Debugf("closing connection %s", c.name)
	if err := c.connection.Close(); err != nil {
		log.Warningf("error closing connection. %v", err)
	}
	c.connected = false
	log.Debugf("connection closed for %s", c.world.Name)
}

// closeFIFO closes the FIFO for the world.
func (c *connection) closeFIFO() {
	name := c.fifo.Name()
	log.Debugf("closing and deleting FIFO %s", name)
	if err := c.fifo.Close(); err != nil {
		log.Warningf("error closing FIFO for reading %s. %v", c.world.Name, err)
	}
	if err := syscall.Unlink(name); err != nil {
		log.Warningf("error unlinking FIFO for %s. %v", c.world.Name, err)
	}
	log.Debugf("FIFO %s closed and deleted for %s", name, c.world.Name)
}

// closeOutputs closes open outfiles.
func (c *connection) closeOutputs() {
	for _, out := range c.outputs {
		log.Debugf("closing output file %s for %s", out.name, c.world.Name)
		if err := out.output.Close(); err != nil {
			log.Warningf("error closing output %s for %s. %v", out.name, c.world.Name, err)
		}
		log.Debugf("output file %s for %s closed", out.name, c.world.Name)
		if !out.global {
			continue
		}
		if c.world.Log {
			rotateTo := c.getLogFile(fmt.Sprintf("%s.log", c.getTimestamp()))
			if err := os.Rename(out.name, rotateTo); err != nil {
				log.Warningf("unable to rotate log file %s, you'll need to do that on your own. %v", out.name, err)
				continue
			}
			log.Debugf("output file for %s rotated", c.world.Name)
		} else {
			if err := os.Remove(out.name); err != nil {
				log.Warningf("unable to remove outfile %s", out.name)
			}
		}
	}
}

// removeWorkingDir removes the (hopefully empty) working directory for the
// connection.
func (c *connection) removeWorkingDir() {
	workingDir := c.getConnectionFile("")
	log.Debugf("removing working directory %s", workingDir)
	if err := os.Remove(workingDir); err != nil {
		log.Errorf("unable to remove working directory %s: %v", workingDir, err)
	}
}

// cleanup cleans up the connection's environment on disk.
func (c *connection) cleanup() {
	log.Debugf("cleaning up connection's environment on disk for %s", c.name)
	c.closeFIFO()
	c.closeOutputs()
	c.removeWorkingDir()
}

// Write sends data to the connection via the FIFO file
func (c *connection) Write(in []byte) (int, error) {
	fname := c.getConnectionFile(inFile)
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_APPEND, os.ModeNamedPipe)
	if err != nil {
		log.Warningf("could not open FIFO for %s! %v", c.world.Name, err)
		return 0, err
	}
	defer f.Close()
	out, err := fmt.Fprintln(f, string(in))
	if err != nil {
		return 0, err
	}
	return out, nil
}

// Close closes the connection and all open files.
func (c *connection) Close() error {
	if !c.connected {
		log.Debugf("%s already closed", c.name)
		return nil
	}
	log.Tracef("closing connection %s", c.name)
	c.disconnect <- true
	if <-c.disconnected {
		c.closeConnection()
		c.cleanup()

		log.Infof("quit %s at %s", c.world.Name, c.getTimestamp())
	}
	return nil
}

// Open opens the connection and all output files.
func (c *connection) Open() error {
	log.Debugf("connecting to %s", c.name)
	var err error

	log.Tracef("creating FIFO for %s", c.name)
	if err = c.makeFIFO(); err != nil {
		return err
	}

	log.Tracef("creating outfile for %s", c.name)
	name := c.getConnectionFile(outFile)
	globalOut := &output{
		name:   name,
		global: true,
		output: nil,
	}
	if err = c.makeLogfile(globalOut); err != nil {
		log.Criticalf("could not create output file for %s: %v", c.name, err)
		c.cleanup()
		return err
	}
	c.outputs = append(c.outputs, globalOut)

	if err = c.connect(); err != nil {
		log.Errorf("could not connect to %s! %v", c.name, err)
		c.cleanup()
		return err
	}
	log.Infof("connected to %s at %s", c.name, c.getTimestamp())

	c.disconnect = make(chan bool)
	c.disconnected = make(chan bool)
	go c.readToFile()
	go c.readToConn()

	return nil
}

// GetConnectionName gets the name of the connection (the connectStr, usually).
func (c *connection) GetConnectionName() string {
	return c.name
}

// GetDisplayName gets the world's display name.
func (c *connection) GetDisplayName() string {
	return c.world.DisplayName
}

// AddOutput creates an output struct with the given io.WriteCloser. This can
// be a file, of course, but many other things as well, including the buffer
// that the UI uses.
func (c *connection) AddOutput(name string, w io.WriteCloser) {
	log.Debugf("creating output %s for %s", name, c.name)
	c.outputs = append(c.outputs, &output{
		name:   name,
		global: false,
		output: w,
	})
}

// NewConnection creates a new conneciton with the given world. One can
// also specify whether or not to use SSL, allow insecure SSL certs, and
// whether to log all output by default.
func NewConnection(name string, w config.World, s config.Server, cfg *config.Config) (*connection, error) {
	log.Debugf("creating a new connection %s for world %s", name, w.Name)
	c := &connection{
		name:      name,
		world:     w,
		server:    s,
		config:    cfg,
		connected: false,
	}

	log.Tracef("ensuring connection working directory")
	if err := os.MkdirAll(c.getConnectionFile(""), 0755); err != nil {
		log.Criticalf("unable to ensure connection directory! %v", err)
		return nil, err
	}

	log.Tracef("ensuring world log directory")
	if err := os.MkdirAll(c.getLogFile(""), 0755); err != nil {
		log.Criticalf("unable to ensure log directory! %v", err)
		return nil, err
	}

	// Look up hostname early on as a network connectivity check.
	log.Tracef("looking up hostname")
	if err := c.lookupHostname(); err != nil {
		log.Errorf("could not look up hostname %s for %s", c.server.Host, c.name)
		return nil, err
	}

	return c, nil
}
