package main

import (
	"flag"
	"log"

	zmq "github.com/pebbe/zmq4"
)

var (
	backendAddress  string
	frontendAddress string
)

func main() {
	flag.StringVar(&backendAddress, "backend", "", "Backend address of the proxy")
	flag.StringVar(&frontendAddress, "frontend", "", "Frontend address of the proxy")
	flag.Parse()

	if backendAddress == "" || frontendAddress == "" {
		flag.PrintDefaults()
		log.Panicln("Command line is incorrect")
	}

	context, err := zmq.NewContext()
	if err != nil {
		log.Panicln(err)
	}
	defer context.Term()

	frontend, err := context.NewSocket(zmq.XSUB)
	if err != nil {
		log.Panicln(err)
	}
	defer frontend.Close()

	if err := frontend.Bind(frontendAddress); err != nil {
		log.Panicln(err)
	}

	backend, err := context.NewSocket(zmq.XPUB)
	if err != nil {
		log.Panicln(err)
	}
	defer backend.Close()

	if err := backend.Bind(backendAddress); err != nil {
		log.Panicln(err)
	}

	if err := zmq.Proxy(frontend, backend, nil); err != nil {
		log.Panicln(err)
	}
}
