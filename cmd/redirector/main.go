package main

import (
	"io"
	"net"

	"github.com/pkg/errors"
)

func main() {
	r := NewRedirect(":8080", "ya.ru:80")
	r.listen()
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
