package fifo

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/juju/loggo"
)

const READ_DELAY = 100 * time.Millisecond

var log = loggo.GetLogger("stimmtausch.connection")

// MakeFIFO creates the FIFO file for the world, used to manage the information
// sent to and recieved from the connection.
func MakeFIFO(file string, readonly bool) (*os.File, error) {
	log.Tracef("creating FIFO file %s", file)
	var err error

	log.Tracef("checking if fifo exists")
	if _, err = os.Stat(file); err == nil {
		log.Errorf("fifo %s already exists!", file)
		return nil, fmt.Errorf("fifo %s already exists, cowardly not continuing", file)
	}

	log.Tracef("making FIFO")
	if err = syscall.Mkfifo(file, 0644); err != nil {
		log.Errorf("unable to make FIFO %s!", file, err)
		return nil, err
	}
	log.Tracef("FIFO created as %s", file)

	log.Tracef("opening FIFO")
	var fifo *os.File

	// Build flags: read-only if we've been asked for such.
	flags := syscall.O_NONBLOCK
	if readonly {
		flags = flags | os.O_RDONLY
	}

	if fifo, err = os.OpenFile(file, flags, os.ModeNamedPipe); err != nil {
		log.Errorf("unable to open FIFO %s! %v", file, err)
		return nil, err
	}
	log.Debugf("FIFO opened as %s", fifo.Name())
	return fifo, nil
}

// OpenFIFO attempts to open a named pipe that already exists.
func OpenFIFO(file string, readonly bool) (*os.File, error) {
	log.Tracef("opening FIFO file %s", file)
	var err error

	log.Tracef("checking if fifo exists")
	if _, err = os.Stat(file); err != nil {
		log.Errorf("fifo %s does not exist!", file)
		return nil, fmt.Errorf("fifo %s does not exist, cowardly not continuing", file)
	}

	flags := syscall.O_NONBLOCK
	if readonly {
		flags = flags | os.O_RDONLY
	}

	var fifo *os.File
	if fifo, err = os.OpenFile(file, flags, os.ModeNamedPipe); err != nil {
		log.Errorf("unable to open FIFO %s! %v", file, err)
		return nil, err
	}
	log.Debugf("FIFO opened as %s", fifo.Name())
	return fifo, nil
}

// OpenOrMakeFIFO attempts to open a named pipe. If it fails because the FIFO
// doesn't exist, it tries to create it instead.
func OpenOrMakeFIFO(file string, readonly bool) (*os.File, error) {
	log.Tracef("trying to open FIFO file %s first", file)
	f, err := OpenFIFO(file, readonly)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return nil, err
	} else if err == nil {
		return f, nil
	}

	log.Tracef("trying to create FIFO file %s instead", file)
	return MakeFIFO(file, readonly)
}

func CloseFIFO(fifo *os.File, unlink bool) error {
	if err := fifo.Close(); err != nil {
		log.Warningf("error closing FIFO %s. %v", f.Name(), err)
		return err
	}
	return nil
}

func UnlinkFIFO(file string) error {
	if unlink {
		if err := syscall.Unlink(file); err != nil {
			log.Warningf("error unlinking FIFO %s. %v", file, err)
			return err
		}
	}
}
