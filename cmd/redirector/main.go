package main

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/pkg/errors"
)

const (
	listenRetryTimeout = time.Second
)

func main() {
	conf, err := parseConfig()
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	for _, redirect := range conf.Redirects {
		r := NewRedirect(redirect.Src, redirect.Dst)
		go r.listenWithRetry()
	}

	select {} // eternal sleep
}

func NewRedirect(source, dest string) *Redirect {
	return &Redirect{
		source: source,
		dest:   dest,
	}
}

type Redirect struct {
	source string
	dest   string
}

func (r *Redirect) listenWithRetry() {
	ticker := time.NewTicker(listenRetryTimeout)
	for {
		log.Printf("Started listening %s -> %s", r.source, r.dest)
		err := r.listen()
		log.Printf("Listen %s -> %s ended, retrying... err=%v", r.source, r.dest, err)
		<-ticker.C
	}
}

func (r *Redirect) listen() error {
	listener, err := net.Listen("tcp", r.source)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return errors.Wrap(err, "listener.Accept")
		}
		go r.handleConnection(conn)
	}
}

func (r *Redirect) handleConnection(incoming net.Conn) error {
	defer incoming.Close()

	outgoing, err := net.Dial("tcp", r.dest)
	if err != nil {
		errors.Wrap(err, "net.Dial")
	}
	defer outgoing.Close()

	go func() {
		io.Copy(incoming, outgoing)
		incoming.Close()
		outgoing.Close()
	}()
	_, err = io.Copy(outgoing, incoming)
	return err
}
