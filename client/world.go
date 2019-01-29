package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	"github.com/juju/loggo"
)

var log = loggo.GetLogger("stimmtausch.client.connection")

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
	output interface {
		Write([]byte) (int, error)
		Close() error
	}
}

// conn stores all connection settings
type world struct {
	name       string
	host       string
	port       uint
	addr       *net.TCPAddr
	ssl        bool
	insecure   bool
	workingDir string
	connection net.Conn
	log        bool
	fifo       *os.File
	outputs    []*output
	disconnect chan bool
}

// getWorldFile returns a file (or directory) name within the scope of the
// world. These live in $HOME/.config/stimmtausch/{worldname}.
func (w *world) getWorldFile(name string) (string, error) {
	u, err := user.Current()
	if err != nil {
		log.Criticalf("cannot get current user! %v", err)
		return "", err
	}
	homeDir := u.HomeDir
	filename := filepath.Join(homeDir, ".config", "stimmtausch", w.name, name)
	log.Tracef("file path: %s", filename)
	return filename, nil
}

// lookupHostname gets the TCP address for the world's hostname.
func (w *world) lookupHostname() error {
	addr, err := net.ResolveTCPAddr("tcp", w.host)
	if err != nil {
		log.Errorf("unable to resolve host %s: %v", w.host, err)
		return err
	}
	w.addr = addr
	log.Debugf("host %s resolves to %v", w.host, w.addr)
	return nil
}

// makeFIFO creates the FIFO file for the world, used to manage the information
// sent to and recieved from the connection.
func (w *world) makeFIFO() error {
	file, err := w.getWorldFile("_fifo")
	if err != nil {
		return err
	}
	if _, err = os.Stat(file); err == nil {
		log.Criticalf("FIFO for connection %s already exists!", w.name)
	}
	if err = syscall.Mkfifo(file, 0644); err != nil {
		log.Criticalf("unable to make FIFO for %s!", w.name, err)
		return err
	}
	log.Debugf("FIFO created as %s", file)
	if w.fifo, err = os.OpenFile(file, os.O_RDONLY|syscall.O_NONBLOCK, os.ModeNamedPipe); err != nil {
		log.Criticalf("unable to open FIFO for %s! %v", file)
		return err
	}
	log.Debugf("FIFO opened as", w.fifo.Name())
	return nil
}

// makeLogfile creates a logfile from a given name.
func (w *world) makeLogfile(out *output) error {
	_, err := os.Stat(out.name)
	if err == nil {
		fmt.Printf("Warning: %v already exists; appending.\n", out.name)
	}
	out.output, err = os.OpenFile(out.name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Warningf("could not open logfile %s for %s, not logging. %v", out.name, w.name, err)
		return err
	}
	log.Debugf("logfile created as %s", out.name)
	return nil
}

// connect creates a TCP connection to the world's TCP address. It connects
// over SSL if available and, if requested, will use insecure certs.
func (w *world) connect() error {
	var err error
	if w.ssl {
		var conf *tls.Config
		if w.insecure {
			conf = &tls.Config{InsecureSkipVerify: true}
		} else {
			conf = &tls.Config{ServerName: w.host}
		}
		if w.connection, err = tls.Dial("tcp", w.addr.String(), conf); err != nil {
			log.Criticalf("unable to dial %v over SSL for %s! %v", w.addr, w.name, err)
			return err
		}
		log.Debugf("connected to server over SSL for %s", w.name)
	} else {
		if w.connection, err = net.DialTCP("tcp", nil, w.addr); err != nil {
			log.Criticalf("unable to dial %v for %s! %v", w.addr, w.name, err)
			return err
		}
		log.Debugf("connected to server for %s", w.name)
	}
	//	// We keep alive for mucks
	//	if err = w.connection.SetKeepAlive(true); err != nil {
	//		log.Warningf("unable to set keep alive for %s - you may get booted. %v", w.name, err)
	//	}
	//	if err = w.connection.SetKeepAlivePeriod(keepalive); err != nil {
	//		log.Warningf("unable to set keep alive period for %s - you may get booted. %v", w.name, err)
	//	}
	return nil
}

// readToConn reads from the FIFO and sends to the connection.
func (w *world) readToConn() {
	tmpError := fmt.Sprintf("read %v: resource temporarily unavailable", w.fifo.Name())
	for {
		select {
		case <-w.disconnect:
			log.Debugf("%s received disconnect; returning", w.name)
			return
		default:
			// This pause between reads from the FIFO is the difference between 0.2%
			// and 100% cpu usage when idle. Also without this you will get excessive
			// "read %v: resource temporarily unavailable" errors on some OSes.
			time.Sleep(fifoReadDelay)
			buf := make([]byte, bufferSize)
			bytesIn, err := w.fifo.Read(buf)
			if err != nil && err.Error() != "EOF" && err.Error() != tmpError {
				log.Criticalf("FIFO broke??¿? World %s. %v", w.name, err)
			} else if bytesIn == 0 {
				continue
			}
			log.Tracef("%d bytes read from FIFO", bytesIn)
			bytesOut, err := w.connection.Write(buf[:bytesIn])
			if err != nil {
				log.Criticalf("FIFO broke??¿? World %s. %v", w.name, err)
			}
			log.Tracef("%d bytes written to file", bytesOut)
		}
	}
}

// readToFile reads from the FIFO and writes to outfiles.
func (w *world) readToFile() {
	for {
		for _, out := range w.outputs {
			reader := bufio.NewReader(w.connection)
			tp := textproto.NewReader(reader)
			line, err := tp.ReadLine()
			if err != nil {
				log.Warningf("server disconnected with %v", err)
				if _, err := fmt.Fprintln(out.output, "\n~Connection lost at %v\n", getTimestamp()); err != nil {
					log.Warningf("unable to write to output %s for %s. %v", out.name, w.name, err)
				}
				w.disconnect <- true
				return
			}
			log.Tracef("%d characters read from %s", len(line), w.name)

			bytesOut, err := fmt.Fprintln(out.output, line)
			if err != nil {
				log.Warningf("unable to write to output %s for world %s. %v", out.name, w.name, err)
			}
			log.Tracef("%d bytes written to output for %s", bytesOut, w.name)
		}
	}
}

// closeConnection closes the world's TCP connection.
func (w *world) closeConnection() {
	if err := w.connection.Close(); err != nil {
		log.Warningf("error closing connection. %v", err)
	}
	log.Debugf("connection closed for %s", w.name)
}

// closeFIFO closes the FIFO for the world.
func (w *world) closeFIFO() {
	name := w.fifo.Name()
	log.Debugf("closing and deleting FIFO %s", name)
	if err := w.fifo.Close(); err != nil {
		log.Warningf("error closing FIFO for %s. %v", w.name, err)
	}
	if err := syscall.Unlink(name); err != nil {
		log.Warningf("error unlinking FIFO for %s. %v", w.name, err)
	}
	log.Debugf("FIFO %s closed and deleted for %s", name, w.name)
}

// closeOut closes open outfiles.
func (w *world) closeOut() {
	for _, out := range w.outputs {
		log.Debugf("closing output file %s for %s", out.name, w.name)
		if err := out.output.Close(); err != nil {
			log.Warningf("error closing output %s for %s. %v", out.name, w.name, err)
		}
		log.Debugf("output file %s for %s closed", out.name, w.name)
		if !out.global {
			continue
		}
		if err := os.Rename(outFile, fmt.Sprintf("%s.log", getTimestamp())); err != nil {
			log.Debugf("error rotating output file %s for %s, will keep using the same file. %v", err)
		}
		log.Debugf("output file for %s rotated", w.name)
	}
}

// Close closes the connection and all open files.
func (w *world) Close() {
	w.disconnect <- true

	w.closeConnection()
	w.closeFIFO()
	w.closeOut()

	log.Infof("quit %s at %s", w.name, getTimestamp())
}

// Open opens the connection and all output files.
func (w *world) Open() error {
	log.Infof("connected to %s at %s", w.name, getTimestamp())

	// Make the in FIFO
	if err := w.makeFIFO(); err != nil {
		return err
	}

	// Make the default log file if requested
	if w.log {
		name, err := w.getWorldFile(fmt.Sprintf("%s.log", getTimestamp()))
		if err != nil {
			log.Criticalf("requested logging for %s but can't comply! %v", w.name, err)
			w.closeFIFO()
			return err
		}
		logfile := &output{
			name:   name,
			global: true,
			output: nil,
		}
		if err = w.makeLogfile(logfile); err != nil {
			return err
		}
		w.outputs = append(w.outputs, logfile)
	}

	w.connect()

	w.disconnect = make(chan bool)
	go w.readToFile()
	go w.readToConn()

	return nil
}

// NewWorld creates a new world with the given name, host, and port. One can
// also specify whether or not to use SSL, allow insecure SSL certs, and
// whether to log all output by default.
func NewWorld(name, host string, port uint, ssl, insecure, logOutput bool) (*world, error) {
	w := &world{
		name:     name,
		host:     host,
		port:     port,
		ssl:      ssl,
		insecure: insecure,
		log:      logOutput,
	}
	var err error
	if w.workingDir, err = w.getWorldFile(""); err != nil {
		log.Criticalf("unable to ensure config directory! %v", err)
		panic(err)
	}
	if err = os.MkdirAll(w.workingDir, 0755); err != nil {
		log.Criticalf("unable to ensure config directory! %v", err)
		panic(err)
	}
	if err := w.lookupHostname(); err != nil {
		return nil, err
	}

	return w, nil
}
