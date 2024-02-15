// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package connection

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/juju/loggo"

	"github.com/makyo/stimmtausch/config"
	"github.com/makyo/stimmtausch/signal"
	"github.com/makyo/stimmtausch/util"
)

var (
	log    = loggo.GetLogger("stimmtausch.connection")
	userRe = regexp.MustCompile("\\$username")
	passRe = regexp.MustCompile("\\$password")
	zwnjRe = regexp.MustCompile("\u200c$")
)

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
func (c *Connection) getTimestamp() string {
	return time.Now().Format(c.config.Client.Logging.TimeString)
}

// conn stores all connection settings
type Connection struct {
	// The name (usually connectStr) of the connection.
	name string

	// The world and server to which this connection belongs.
	// These are maintained separately from the app config as they may be
	// connected and passed in for settings not in the user's config.
	world  config.World
	server config.Server

	// The app config.
	config *config.Config

	// The signal dispatcher.
	env *signal.Dispatcher

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

	// A channel to listen for signal events.
	listener chan signal.Signal

	// Whether or not the server is connected.
	Connected bool
}

// lookupHostname gets the TCP address for the world's hostname.
func (c *Connection) lookupHostname() error {
	log.Tracef("attempting to resolve %s:%d", c.server.Host, c.server.Port)
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
func (c *Connection) getConnectionFile(name string) string {
	return filepath.Join(c.config.WorkingDir, c.name, name)
}

// getLogFile returns a file (or directory) name within the scope of the
// connection for the sake of logging. These live in
// $HOME/.local/log/stimmtausch/{worldname}.
func (c *Connection) getLogFile(name string) string {
	return filepath.Join(c.config.LogDir, c.name, name)
}

// makeFIFO creates the FIFO file for the world, used to manage the information
// sent to and recieved from the connection.
func (c *Connection) makeFIFO() error {
	log.Tracef("creating FIFO file for %s", c.name)
	file := c.getConnectionFile(inFile)
	var err error

	log.Tracef("checking if FIFO exists")
	if _, err = os.Stat(file); err == nil {
		log.Errorf("FIFO for connection %s already exists!", c.name)
		return fmt.Errorf("FIFO for connection %s already exists, cowardly not continuing", c.name)
	}

	log.Tracef("making FIFO")
	if err = syscall.Mkfifo(file, 0644); err != nil {
		log.Errorf("unable to make FIFO for %s!", c.name, err)
		return err
	}
	log.Tracef("FIFO created as %s", file)

	log.Tracef("opening FIFO")
	if c.fifo, err = os.OpenFile(file, os.O_RDONLY|syscall.O_NONBLOCK, os.ModeNamedPipe); err != nil {
		log.Errorf("unable to open FIFO for reading %s! %v", file, err)
		return err
	}
	log.Debugf("FIFO opened as %s", c.fifo.Name())
	return nil
}

// connect creates a TCP connection to the world's TCP address. It connects
// over SSL if available and, if requested, will use insecure certs.
func (c *Connection) connect() error {
	log.Tracef("creating TCP connection for %s", c.name)
	var err error
	var conn *net.TCPConn

	log.Tracef("creating regular TCP connection")
	if conn, err = net.DialTCP("tcp", nil, c.addr); err != nil {
		log.Errorf("unable to dial %v for %s! %v", c.addr, c.name, err)
		return err
	}

	log.Tracef("attempting to set a keepalive for %s", c.name)
	if err = conn.SetKeepAlive(true); err != nil {
		log.Warningf("unable to set keep alive for %s - you may get booted. %v", c.name, err)
	}
	if err = conn.SetKeepAlivePeriod(keepalive); err != nil {
		log.Warningf("unable to set keep alive period for %s - you may get booted. %v", c.name, err)
	}
	c.connection = conn
	fmt.Fprintln(c.connection, "\xff\xfdCHARSET unicode")
	log.Debugf("connected to server for %s", c.name)

	if c.server.SSL {
		log.Tracef("creating SSL connection")
		var conf *tls.Config
		if c.server.Insecure {
			conf = &tls.Config{InsecureSkipVerify: true}
		} else {
			conf = &tls.Config{ServerName: c.server.Host}
		}
		c.connection = tls.Client(conn, conf)
		log.Debugf("connected to server over SSL for %s", c.name)
	}

	c.Connected = true
	return nil
}

// readToConn reads from the FIFO and sends to the connection.
func (c *Connection) readToConn() {
	log.Tracef("reading from FIFO to connection %s", c.name)
	tmpError := fmt.Sprintf("read %v: resource temporarily unavailable", c.fifo.Name())
	for {
		select {
		case <-c.disconnect:
			log.Debugf("%s received disconnect; returning", c.name)
			c.disconnected <- true
			return
		default:
			// This pause between reads from the FIFO is the difference between 0.2%
			// and 100% cpu usage when idle. Also without this you will get excessive
			// "read %v: resource temporarily unavailable" errors on some OSes.
			time.Sleep(fifoReadDelay)
			scanner := bufio.NewScanner(c.fifo)
			if !scanner.Scan() {
				continue
			}
			text := scanner.Text()
			if len(text) == 0 {
				log.Infof("got an empty string from the buffer, which is weird.")
				continue
			}
			if text[0] == '/' {
				s := strings.SplitN(text[1:], " ", 2)
				if len(s) == 1 {
					s = append(s, "")
				}
				go c.env.Dispatch(s[0], s[1])
				continue
			}
			if err := scanner.Err(); err != nil && err.Error() != tmpError {
				log.Errorf("FIFO broke??¿? connection %s. %v", c.name, err)
				continue
			}
			fmt.Fprintln(c.connection, text)
		}
	}
}

// readToFile reads from the connection and writes to outfiles.
func (c *Connection) readToFile() {
	log.Tracef("reading from connection %s to file", c.name)
	reader := bufio.NewReader(c.connection)
	tp := textproto.NewReader(reader)
	for {
		bareLine, err := tp.ReadLine()
		line := strings.ToValidUTF8(bareLine, "")
		if err != nil {
			if !c.Connected {
				return
			}
			log.Warningf("server disconnected with %v", err)
			disconnectMsg := fmt.Sprintf("\n~Connection lost at %v\n", c.getTimestamp())
			for _, out := range c.outputs {
				if _, err := fmt.Fprintln(out.output, disconnectMsg); err != nil {
					log.Warningf("unable to write to output %s for %s. %v", out.name, c.name, err)
				}
			}
			c.Close()
			c.Connected = false
			return
		}
		log.Tracef("%d characters read from %s", len(line), c.name)

		log.Tracef("running triggers against line")
		var errs, triggerErrs []error
		var applies, gag, logAnyway bool
		orig := line
		for _, trigger := range c.config.CompiledTriggers {
			applies, line, triggerErrs = trigger.Run(c.world.Name, line, c.config)
			if len(triggerErrs) != 0 {
				errs = append(errs, triggerErrs...)
			}
			if applies && trigger.Type == "gag" {
				log.Tracef("gag %+v applies", trigger)
				gag = true
				logAnyway = trigger.LogAnyway
			}
		}
		// Some worlds end a line with a ZWNJ (\u200c) in order to aid in triggers in wrapped text. Remove before printing
		line = zwnjRe.ReplaceAllString(line, "")
		if len(errs) != 0 {
			log.Errorf("errors encountered processing triggers: %q", errs)
		}
		for _, out := range c.outputs {
			if gag && !(logAnyway && out.global) {
				continue
			}
			toWrite := orig
			if out.supportsANSI {
				toWrite = line
			}
			bytesOut, err := fmt.Fprintln(out.output, toWrite)
			if err != nil {
				log.Warningf("unable to write to output %s for connection %s. %v", out.name, c.name, err)
			}
			log.Tracef("%d bytes written to output %s for %s", bytesOut, out.name, c.name)
		}
	}
}

// closeConnection closes the world's TCP connection.
func (c *Connection) closeConnection() {
	if !c.Connected {
		log.Debugf("%s already closed", c.name)
		return
	}
	log.Tracef("closing connection %s", c.name)
	if err := c.connection.Close(); err != nil {
		log.Warningf("error closing connection. %v", err)
	}
	c.Connected = false
	log.Debugf("connection closed for %s", c.name)
}

// closeFIFO closes the FIFO for the world.
func (c *Connection) closeFIFO() {
	name := c.fifo.Name()
	log.Tracef("closing and deleting FIFO %s", name)
	if err := c.fifo.Close(); err != nil {
		log.Warningf("error closing FIFO for reading %s. %v", c.name, err)
	}
	if err := syscall.Unlink(name); err != nil {
		log.Warningf("error unlinking FIFO for %s. %v", c.name, err)
	}
	log.Debugf("FIFO %s closed and deleted for %s", name, c.name)
}

// removeWorkingDir removes the (hopefully empty) working directory for the
// connection.
func (c *Connection) removeWorkingDir() {
	workingDir := c.getConnectionFile("")
	log.Tracef("removing working directory %s", workingDir)
	if err := os.Remove(workingDir); err != nil {
		log.Errorf("unable to remove working directory %s: %v", workingDir, err)
	}
	log.Debugf("working directory %s removed for %s", workingDir, c.name)
}

// cleanup cleans up the connection's environment on disk.
func (c *Connection) cleanup() {
	log.Tracef("cleaning up connection's environment on disk for %s", c.name)
	c.closeFIFO()
	c.closeOutputs()
	c.removeWorkingDir()
}

// Write sends data to the connection via the FIFO file
func (c *Connection) Write(in []byte) (int, error) {
	fname := c.getConnectionFile(inFile)
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_APPEND, os.ModeNamedPipe)
	if err != nil {
		log.Warningf("could not open FIFO for %s! %v", c.name, err)
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
func (c *Connection) Close() error {
	if !c.Connected {
		log.Debugf("%s already closed", c.name)
		return nil
	}
	log.Tracef("closing connection %s", c.name)
	c.disconnect <- true
	if <-c.disconnected {
		c.closeConnection()
		c.cleanup()
		c.env.Dispatch("_client:disconnected", c.name)

		log.Infof("quit %s at %s", c.name, c.getTimestamp())
	}
	return nil
}

// listen listens for events from the signal environment, then does nothing (but
// does it splendidly)
func (c *Connection) listen() {
	for {
		res := <-c.listener
		switch res.Name {
		case "log":
			if err := c.parseLogSignal(res.Payload); err != nil {
				log.Errorf("error executing log command: %v", err)
			}
		default:
			continue
		}
	}
}

// Open opens the connection and all output files.
func (c *Connection) Open() error {
	log.Tracef("connecting to %s", c.name)
	var err error

	log.Tracef("creating FIFO for %s", c.name)
	if err = c.makeFIFO(); err != nil {
		return err
	}

	log.Tracef("creating outfile for %s", c.name)
	name := c.getConnectionFile(outFile)
	globalOut := &output{
		name:         name,
		global:       true,
		output:       nil,
		supportsANSI: true, // Global out supports ANSI, which is stripped during rotation.
	}
	if err = c.makeLogfile(globalOut); err != nil {
		log.Errorf("could not create output file for %s: %v", c.name, err)
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

	log.Tracef("listening for signals")
	c.listener = make(chan signal.Signal)
	go c.listen()
	c.env.AddListener("connection", c.listener)

	c.disconnect = make(chan bool)
	c.disconnected = make(chan bool)
	go c.readToFile()
	go c.readToConn()

	st, ok := c.config.ServerTypes[c.server.ServerType]
	if ok && c.world.Username != "" && c.world.Password != "" {
		connectStr := st.ConnectString
		connectStr = userRe.ReplaceAllString(connectStr, c.world.Username)
		connectStr = passRe.ReplaceAllString(connectStr, c.world.Password)
		fmt.Fprintln(c, connectStr)
	}

	return nil
}

// GetConnectionName gets the name of the connection (the connectStr, usually).
func (c *Connection) GetConnectionName() string {
	return c.name
}

// GetDisplayName gets the world's display name.
func (c *Connection) GetDisplayName() string {
	return c.world.DisplayName
}

// GetMaxBuffer returns the max buffer lengh for a server.
func (c *Connection) GetMaxBuffer() uint {
	return c.server.MaxBuffer
}

// NewConnection creates a new conneciton with the given world. One can
// also specify whether or not to use SSL, allow insecure SSL certs, and
// whether to log all output by default.
func NewConnection(name string, w config.World, s config.Server, cfg *config.Config, env *signal.Dispatcher) (*Connection, error) {
	log.Tracef("creating a new connection %s for world %s", name, w.Name)
	c := &Connection{
		name:      name,
		world:     w,
		server:    s,
		config:    cfg,
		env:       env,
		Connected: false,
	}

	log.Tracef("ensuring connection working directory")
	if err := util.EnsureDir(c.getConnectionFile("")); err != nil {
		log.Errorf("unable to ensure connection directory! %v", err)
		return nil, err
	}

	log.Tracef("ensuring world log directory")
	if err := util.EnsureDir(c.getLogFile("")); err != nil {
		log.Errorf("unable to ensure log directory! %v", err)
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
