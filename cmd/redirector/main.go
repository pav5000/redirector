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

	LogLevelSilent     = 0
	LogLevelListens    = 1
	LogLevelDialErrors = 2
	LogLevelAllConns   = 3
)

func main() {
	conf, err := parseConfig()
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	for _, redirect := range conf.Redirects {
		r := NewRedirect(redirect.Src, redirect.Dst, conf.Verbose)
		go r.listenWithRetry()
	}

	if conf.Verbose == 0 {
		log.Println("Started")
	}

	select {} // eternal sleep
}

func NewRedirect(source, dest string, verbose int) *Redirect {
	return &Redirect{
		verbose: verbose,
		source:  source,
		dest:    dest,
	}
}

type Redirect struct {
	verbose int
	source  string
	dest    string
}

func (r *Redirect) listenWithRetry() {
	ticker := time.NewTicker(listenRetryTimeout)
	for {
		r.logPrintf(LogLevelListens, "Started listening %s -> %s", r.source, r.dest)
		err := r.listen()
		r.logPrintf(LogLevelListens, "Listen %s -> %s ended, retrying... err=%s", r.source, r.dest, err.Error())
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
		go func() {
			err := r.handleConnection(conn)
			if err != nil {
				r.logPrintf(LogLevelAllConns, "Connection %s -> %s closed with error %s", conn.RemoteAddr().String(), r.dest, err.Error())
			} else {
				r.logPrintf(LogLevelAllConns, "Connection %s -> %s closed", conn.RemoteAddr().String(), r.dest)
			}
		}()
	}
}

func (r *Redirect) handleConnection(incoming net.Conn) error {
	defer incoming.Close()

	r.logPrintf(LogLevelAllConns, "New connection %s -> %s", incoming.RemoteAddr().String(), r.dest)

	outgoing, err := net.Dial("tcp", r.dest)
	if err != nil {
		r.logPrintf(LogLevelDialErrors, "Cannot open new connection to %s: %s", r.dest, err.Error())
		return errors.Wrap(err, "net.Dial")
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

func (r *Redirect) logPrintf(verbose int, format string, v ...interface{}) {
	if r.verbose < verbose {
		return
	}
	log.Printf(format, v...)
}
