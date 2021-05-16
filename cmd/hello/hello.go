package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"github.com/nats-io/nats.go"
)

func main() {

	parser := argparse.NewParser("flags", "Ping Pong tool for NATS connectivity")
	subscriberFlag := parser.Flag("v", "verbose", &argparse.Options{Help: "Enable subcriber mode"})
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
	defer nc.Drain()
	defer nc.Close()

	if *subscriberFlag {
		// Simple Async Subscriber
		fmt.Printf("Start reciving")
		nc.Subscribe("pingpong", func(m *nats.Msg) {
			fmt.Printf("Received a message: %s\n", string(m.Data))
		})
		sleep()

	} else {
		rand.Seed(time.Now().UnixNano())
		seed := rand.Intn(300)
		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("Sent a message: %s, seed %d, : count %d/10\n", "Yello", seed, i+1)
			fmt.Print(msg)
			nc.Publish("pingpong", []byte(msg))
		}
	}

}

func sleep() {
	for {
		time.Sleep(1000)
	}
}
