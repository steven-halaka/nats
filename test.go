package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://a:a@localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	fmt.Printf("Connected to %s\n", nc.ConnectedUrl())
	fmt.Printf("Connect urls: %s\n", nc.Servers())

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	stream, err := js.AddStream(&nats.StreamConfig{
		Name:     "tk",
		Subjects: []string{"*"},
		Storage:  nats.FileStorage,
		Replicas: 3,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatal(err)
	}

	_, err = js.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received message on %s: %s\n", m.Subject, string(m.Data))
		m.Ack()
	},
		nats.Durable("durable-foo"),
		nats.ManualAck(),
		nats.AckExplicit(),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 3; i++ {
		msg := fmt.Sprintf("Message %d", i)
		_, err = js.Publish("foo", []byte(msg))
		if err != nil {
			log.Fatal(err)
		}
	}

	time.Sleep(1 * time.Second)

	state, _ := json.MarshalIndent(stream.State, "", "  ")
	fmt.Println(string(state))
}
