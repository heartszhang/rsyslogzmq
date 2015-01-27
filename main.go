package main

import (
	"bufio"
	"flag"
	"log"
	"net"
)

var option = struct {
	sock     string
	zmq_addr string
	verbose  bool
}{
	sock:     "localhost:4514",
	zmq_addr: "ipc://rsyslog.zmq.pull",
}

func init() {
	flag.StringVar(&option.sock, "sock", option.sock, "rsyslog upstream socket")
	flag.BoolVar(&option.verbose, "verbose", option.verbose, "verbose mode")
}

func main() {
	flag.Parse()
	ln, err := net.Listen("tcp", option.sock)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	for {
		if conn, err := ln.Accept(); err == nil {
			go handle_connection(conn)
		} else {
			log.Fatal(err)
		}
	}
}

func handle_connection(conn net.Conn) {
	log.Println("client-start", conn.RemoteAddr())
	defer conn.Close()
	doc_chan := make(chan []byte, 64)
	defer close(doc_chan)
	go zmq_push(doc_chan, conn)
	var err error
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() && err == nil {
		doc_chan <- cp(scanner.Bytes())
	}
	log.Println("client-close", conn.RemoteAddr(), err)
}

func cp(i []byte) []byte {
	o := make([]byte, len(i))
	copy(o, i)
	return o
}
