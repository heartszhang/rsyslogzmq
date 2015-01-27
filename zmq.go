package main

import (
	"log"
	"net"

	"github.com/pebbe/zmq4"
)

func zmq_push(channel chan []byte, conn net.Conn) {
	defer conn.Close()
	pusher, err := zmq4.NewSocket(zmq4.PUB)
	if err == nil {
		defer pusher.Close()
	}
	if err = pusher.Bind(option.zmq_addr); err != nil {
		log.Println(err)
		return
	}
	// dont block send
	for doc := range channel {
		if _, err := pusher.SendBytes(doc, 0); err != nil {
			log.Println("zmq-push", err)
			break
		}
	}
}
