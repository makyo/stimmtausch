package eventstream

import (
	"os"

	"github.com/makyo/stimmtausch/connection/fifo"
)

type eventStream struct {
	sendfifo         *os.File
	recvfifo         *os.File
	basedir          string
	clientID         int
	handlers         []func(Event)
	listening        bool
	stop             chan bool
	stopped          chan bool
	hostNextClientID int
}

func (es *eventStream) listen() {
	for {
		select {
		case <-es.stop:
			es.listening = false
			es.stopped <- true
			return
		default:
			time.Sleep(fifo.READ_DELAY)
			// read from recvfifo, act on that
		}
	}
}

func (es *eventStream) AddEventHandler(handler func(Event)) {
	es.handlers = append(es.handlers, handler)
}

func (es *eventStream) Init() error {
	if es.listening {
		return nil
	}
	if es.hostNextClientID > 0 {
		sendfifo, err := fifo.OpenOrMakeFIFO(filepath.Join(basedir, "send"), false)
		if err != nil {
			return err
		}
		recvfifo, err := fifo.OpenOrMakeFIFO(filepath.Join(basedir, "recv"), true)
		if err != nil {
			return err
		}
		es.sendfifo = recvfifo
		es.recvfifo = sendfifo
	} else {
		sendfifo, err := fifo.OpenFIFO(filepath.Join(basedir, "send"), false)
		if err != nil {
			return err
		}
		recvfifo, err := fifo.OpenFIFO(filepath.Join(basedir, "recv"), true)
		if err != nil {
			return err
		}
		es.sendfifo = sendfifo
		es.recvfifo = recvfifo
	}
	go es.listen()
	if es.clientID < 1 {
		if err := es.handshake(); err != nil {
			return err
		}
	}
}

func (es *eventStream) Stop() error {
	es.stop <- true
	<-es.stopped
	if err := fifo.Close(es.sendfifo); err != nil {
		return err
	}
	if err := fifo.Close(es.recvfifo); err != nil {
		return err
	}
	es.clientID = 0
	return nil
}

func (es *eventStream) Destroy() error {
	if es.listening {
		return fmt.Errorf("currently listening for events")
	}
}

func NewStream(basedir string) (*eventStream, error) {
	es := &eventStream{
		basedir:          basedir,
		clientID:         0,
		handlers:         []handler{},
		stop:             make(chan bool),
		stopped:          make(chan bool),
		hostNextClientID: 0,
	}
	return es
}

func NewStreamHost(basedir string) (*eventStream, error) {
	es, err := NewStream(basedir, 1)
	if err != nil {
		return nil, err
	}
	es.clientID = 1
	es.hostNextClientID = 2
	es.AddEventHandler(handshakeHandler)
}

func handshakeHandler(event Event) {
}
