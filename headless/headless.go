// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package headless

import (
	"errors"
	"os"
	sig "os/signal"

	"github.com/juju/loggo"

	"github.com/makyo/stimmtausch/client"
	"github.com/makyo/stimmtausch/signal"
)

var (
	log     = loggo.GetLogger("stimmtausch.headless")
	errQuit = errors.New("quitting")
)

type headless struct {
	client         *client.Client
	listener       chan signal.Signal
	signalListener chan os.Signal
	err            chan error
}

func (h *headless) connect(name string) {
	log.Tracef("creating a connection with connection string %s", name)
	conn, ok := h.client.Conn(name)
	if !ok {
		log.Errorf("unable to find connection %s", name)
	}

	log.Tracef("opening connection for %s", name)
	err := conn.Open()
	if err != nil {
		log.Errorf("unable to open connection for %s: %v", name, err)
		h.err <- err
	}
}

func (h *headless) listen() {
	for {
		res := <-h.listener
		switch res.Name {
		case "help":
			log.Infof("Headless Stimmtausch help")
		case "_client:connect":
			h.connect(res.Payload[0])
		default:
			log.Tracef("got unknown signal result %v", res)
		}
	}
}

func (h *headless) mainLoop() error {
	h.err = make(chan error)
	for {
		select {
		case err := <-h.err:
			return err
		case s := <-h.signalListener:
			log.Tracef("received signal to quit %v", s)
			return errQuit
		}
	}
	return nil
}

func (h *headless) Run(done chan bool) {
	h.listener = make(chan signal.Signal)
	go h.listen()
	h.client.Env.AddListener("headless", h.listener)

	if err := h.mainLoop(); err != nil && err != errQuit {
		log.Criticalf("headless unexpectedly quit: %v", err)
	}
	h.client.CloseAll()
	log.Infof("I think we're done?")
	done <- true
}

func New(args []string, client *client.Client) *headless {
	h := &headless{
		client:         client,
		signalListener: make(chan os.Signal, 1),
	}
	sig.Notify(h.signalListener,
		os.Interrupt,
		os.Kill)
	return h
}
