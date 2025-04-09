package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	opts := []nats.Option{
		nats.Name("JetStream Client with Callbacks"),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(-1),
		nats.DisconnectHandler(func(nc *nats.Conn) {
			log.Println("[!] Disconnected from NATS!")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("[+] Reconnected to NATS at %s\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Println("[x] Connection permanently closed.")
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("[!] Error: %v\n", err)
		}),
	}

	// Connect to NATS
	nc, err := nats.Connect("nats://a:a@localhost:4222", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	// JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Add the stream (idempotent)
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "tk",
		Subjects: []string{"*"},
		Storage:  nats.FileStorage,
		Replicas: 3,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatal(err)
	}

	// Subscribe with durable consumer and manual ack
	_, err = js.Subscribe("foo", func(msg *nats.Msg) {
		log.Printf("Received: %s\n", string(msg.Data))
		msg.Ack()
	},
		nats.Durable("durable-foo"),
		nats.ManualAck(),
		nats.AckExplicit(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start background publisher
	go func() {
		count := 1
		for {
			payload := fmt.Sprintf("Ping %d @ %s", count, time.Now().Format(time.RFC3339))
			_, err := js.Publish("foo", []byte(payload))
			if err != nil {
				log.Printf("[!] Publish failed: %v\n", err)
			} else {
				log.Printf("Published: %s\n", payload)
			}
			count++
			time.Sleep(time.Second)
		}
	}()

	// Block forever
	select {}
}
