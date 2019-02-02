package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"os"
	"syscall"
	"time"

	"github.com/juju/loggo"
)

var log = loggo.GetLogger("stimmtausch.client")

// hardcoded program settings
const inFile string = "in"
const outFile string = "out"
const timeString string = "2006-01-02T150405"
const bufferSize int = 1024
const fifoReadDelay = 100 * time.Millisecond
const keepalive = 15 * time.Minute

func getTimestamp() string {
	return time.Now().Format(timeString)
}

type output struct {
	name   string
	global bool
	output io.WriteCloser
}

// conn stores all connection settings
type connection struct {
	name         string
	world        *world
	addr         *net.TCPAddr
	workingDir   string
	connection   net.Conn
	fifo         *os.File
	outputs      []*output
	disconnect   chan bool
	disconnected chan bool
	connected    bool
}

// lookupHostname gets the TCP address for the world's hostname.
func (c *connection) lookupHostname() error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.world.server.host, c.world.server.port))
	if err != nil {
		log.Errorf("unable to resolve host %s: %v", c.world.server.host, err)
		return err
	}
	c.addr = addr
	log.Debugf("host %s resolves to %v", c.world.server.host, c.addr)
	return nil
}

// makeFIFO creates the FIFO file for the world, used to manage the information
// sent to and recieved from the connection.
func (c *connection) makeFIFO() error {
	file, err := c.world.getWorldFile("_fifo")
	if err != nil {
		return err
	}
	if _, err = os.Stat(file); err == nil {
		log.Criticalf("FIFO for connection %s already exists!", c.world.name)
	}
	if err = syscall.Mkfifo(file, 0644); err != nil {
		log.Criticalf("unable to make FIFO for %s!", c.world.name, err)
		return err
	}
	log.Debugf("FIFO created as %s", file)
	if c.fifo, err = os.OpenFile(file, os.O_RDONLY|syscall.O_NONBLOCK, os.ModeNamedPipe); err != nil {
		log.Criticalf("unable to open FIFO for reading %s! %v", file, err)
		return err
	}
	log.Debugf("FIFO opened as", c.fifo.Name())
	return nil
}

// makeLogfile creates a logfile from a given name.
func (c *connection) makeLogfile(out *output) error {
	_, err := os.Stat(out.name)
	if err == nil {
		fmt.Printf("Warning: %v already exists; appending.\n", out.name)
	}
	f, err := os.OpenFile(out.name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Warningf("could not open logfile %s for %s, not logging. %v", out.name, c.world.name, err)
		return err
	}
	out.output = f
	log.Debugf("logfile created as %s", out.name)
	return nil
}

// connect creates a TCP connection to the world's TCP address. It connects
// over SSL if available and, if requested, will use insecure certs.
func (c *connection) connect() error {
	var err error
	if c.world.server.ssl {
		var conf *tls.Config
		if c.world.server.insecure {
			conf = &tls.Config{InsecureSkipVerify: true}
		} else {
			conf = &tls.Config{ServerName: c.world.server.host}
		}
		if c.connection, err = tls.Dial("tcp", c.addr.String(), conf); err != nil {
			log.Criticalf("unable to dial %v over SSL for %s! %v", c.addr, c.world.name, err)
			return err
		}
		log.Debugf("connected to server over SSL for %s", c.world.name)
	} else {
		if c.connection, err = net.DialTCP("tcp", nil, c.addr); err != nil {
			log.Criticalf("unable to dial %v for %s! %v", c.addr, c.world.name, err)
			return err
		}
		log.Debugf("connected to server for %s", c.world.name)
	}
	//	// We keep alive for mucks
	//	if err = c.connection.SetKeepAlive(true); err != nil {
	//		log.Warningf("unable to set keep alive for %s - you may get booted. %v", c.world.name, err)
	//	}
	//	if err = c.connection.SetKeepAlivePeriod(keepalive); err != nil {
	//		log.Warningf("unable to set keep alive period for %s - you may get booted. %v", c.world.name, err)
	//	}
	c.connected = true
	return nil
}

// readToConn reads from the FIFO and sends to the connection.
func (c *connection) readToConn() {
	tmpError := fmt.Sprintf("read %v: resource temporarily unavailable", c.fifo.Name())
	for {
		select {
		case <-c.disconnect:
			log.Debugf("%s received disconnect; returning", c.world.name)
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
				log.Criticalf("FIFO broke??¿? World %s. %v", c.world.name, err)
			} else if bytesIn == 0 {
				continue
			}
			log.Tracef("%d bytes read from FIFO: %s", bytesIn, string(buf[:bytesIn]))
			bytesOut, err := c.connection.Write(buf[:bytesIn])
			if err != nil {
				log.Criticalf("FIFO broke??¿? World %s. %v", c.world.name, err)
			}
			log.Tracef("%d bytes written to connection", bytesOut)
		}
	}
}

// readToFile reads from the FIFO and writes to outfiles.
func (c *connection) readToFile() {
	reader := bufio.NewReader(c.connection)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			log.Warningf("server disconnected with %v", err)
			disconnectMsg := fmt.Sprintf("\n~Connection lost at %v\n", getTimestamp())
			for _, out := range c.outputs {
				if _, err := fmt.Fprintln(out.output, disconnectMsg); err != nil {
					log.Warningf("unable to write to output %s for %s. %v", out.name, c.world.name, err)
				}
			}
			c.disconnect <- true
			return
		}
		log.Tracef("%d characters read from %s", len(line), c.world.name)

		for _, out := range c.outputs {
			bytesOut, err := fmt.Fprintln(out.output, line)
			if err != nil {
				log.Warningf("unable to write to output %s for world %s. %v", out.name, c.world.name, err)
			}
			log.Tracef("%d bytes written to output %s for %s", bytesOut, out.name, c.world.name)
		}
	}
}

// closeConnection closes the world's TCP connection.
func (c *connection) closeConnection() {
	if !c.connected {
		log.Debugf("%s already closed", c.name)
	}
	if err := c.connection.Close(); err != nil {
		log.Warningf("error closing connection. %v", err)
	}
	c.connected = false
	log.Debugf("connection closed for %s", c.world.name)
}

// closeFIFO closes the FIFO for the world.
func (c *connection) closeFIFO() {
	name := c.fifo.Name()
	log.Debugf("closing and deleting FIFO %s", name)
	if err := c.fifo.Close(); err != nil {
		log.Warningf("error closing FIFO for reading %s. %v", c.world.name, err)
	}
	if err := syscall.Unlink(name); err != nil {
		log.Warningf("error unlinking FIFO for %s. %v", c.world.name, err)
	}
	log.Debugf("FIFO %s closed and deleted for %s", name, c.world.name)
}

// closeOut closes open outfiles.
func (c *connection) closeOut() {
	for _, out := range c.outputs {
		log.Debugf("closing output file %s for %s", out.name, c.world.name)
		if err := out.output.Close(); err != nil {
			log.Warningf("error closing output %s for %s. %v", out.name, c.world.name, err)
		}
		log.Debugf("output file %s for %s closed", out.name, c.world.name)
		if !out.global {
			continue
		}
		if err := os.Rename(outFile, fmt.Sprintf("%s.log", getTimestamp())); err != nil {
			log.Debugf("error rotating output file %s for %s, will keep using the same file. %v", err)
		}
		log.Debugf("output file for %s rotated", c.world.name)
	}
}

func (c *connection) Write(in []byte) (int, error) {
	fname, err := c.world.getWorldFile("_fifo")
	if err != nil {
		return 0, err
	}
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_APPEND, os.ModeNamedPipe)
	if err != nil {
		log.Warningf("could not open FIFO for %s! %v", c.world.name, err)
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
		c.closeFIFO()
		c.closeOut()

		log.Infof("quit %s at %s", c.world.name, getTimestamp())
	}
	return nil
}

// Open opens the connection and all output files.
func (c *connection) Open() error {
	log.Infof("connected to %s at %s", c.world.name, getTimestamp())

	// Make the FIFO
	if err := c.makeFIFO(); err != nil {
		return err
	}

	// Make the default log file if requested
	if c.world.log {
		name, err := c.world.getWorldFile(fmt.Sprintf("%s.log", getTimestamp()))
		if err != nil {
			log.Criticalf("requested logging for %s but can't comply! %v", c.world.name, err)
			c.closeFIFO()
			return err
		}
		logfile := &output{
			name:   name,
			global: true,
			output: nil,
		}
		if err = c.makeLogfile(logfile); err != nil {
			return err
		}
		c.outputs = append(c.outputs, logfile)
	}

	c.connect()

	c.disconnect = make(chan bool)
	c.disconnected = make(chan bool)
	go c.readToFile()
	go c.readToConn()

	return nil
}

func (c *connection) World() *world {
	return c.world
}

func (c *connection) Server() *server {
	return c.world.server
}

func (c *connection) GetConnectionName() string {
	return c.name
}

func (c *connection) AddOutput(name string, w io.WriteCloser) {
	c.outputs = append(c.outputs, &output{
		name:   name,
		global: false,
		output: w,
	})
}

// NewConnection creates a new conneciton with the given world. One can
// also specify whether or not to use SSL, allow insecure SSL certs, and
// whether to log all output by default.
func NewConnection(name string, w *world) (*connection, error) {
	c := &connection{
		name:      name,
		world:     w,
		connected: false,
	}
	var err error
	if c.workingDir, err = c.world.getWorldFile(""); err != nil {
		log.Criticalf("unable to ensure config directory! %v", err)
		panic(err)
	}
	if err = os.MkdirAll(c.workingDir, 0755); err != nil {
		log.Criticalf("unable to ensure config directory! %v", err)
		panic(err)
	}
	if err := c.lookupHostname(); err != nil {
		return nil, err
	}

	return c, nil
}
