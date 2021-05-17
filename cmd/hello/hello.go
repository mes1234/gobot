package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"github.com/nats-io/nats.go"
)

type PingPonger struct {
	I string
}

func (pp *PingPonger) Other() string {
	if pp.I == "ping" {
		return "pong"
	} else {
		return "ping"
	}
}

func pingpong(pingpong PingPonger, nc *nats.Conn) {
	log.Printf("Hello I am %v", pingpong.I)
	if pingpong.I == "ping" {
		log.Printf("published to %v", pingpong.I)
		nc.Publish(pingpong.I, []byte(pingpong.I))
	}
	log.Printf("I will start to listen to to %v", pingpong.Other())
	nc.Subscribe(pingpong.Other(), func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
		nc.Publish(pingpong.I, []byte(pingpong.Other()))
		time.Sleep(1 * time.Second)
	})
	sleep()
}

func main() {

	parser := argparse.NewParser("flags", "Ping Pong tool for NATS connectivity")
	pingPong := parser.String("p", "pingpong", &argparse.Options{Help: "Are you ping or pong"})
	serverAddr := parser.String("s", "server", &argparse.Options{Required: true, Help: "server address is requried"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	// Connect to a server
	nc, _ := nats.Connect(*serverAddr)
	// defer nc.Drain()
	// defer nc.Close()

	switch *pingPong {
	case "ping":
		go pingpong(PingPonger{I: "ping"}, nc)
	case "pong":
		go pingpong(PingPonger{I: "pong"}, nc)
	}
	sleep()
}

func sleep() {
	for {
		time.Sleep(1000)
	}
}
