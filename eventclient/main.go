package main

import (
	"fmt"
	protos "github.com/Buzzology/practical-golang/eventserver/protos/eventserver"
	"github.com/nats-io/nats.go"
	natsp "github.com/nats-io/nats.go/encoders/protobuf"
	"os"
)

func main() {

	// Ensure that we've been passed a nats server address
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}

	// Connect to nats
	nc, err := nats.Connect(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Connect to nats (encoded)
	ec, err := nats.NewEncodedConn(nc, natsp.PROTOBUF_ENCODER)
	if err != nil {
		fmt.Errorf("failed to create encoded nats connection: %v", err)
		os.Exit(1)
	}

	defer ec.Close()

	// Receive our standard messages
	ec.Subscribe("Messaging.Text.Standard", func(m *protos.TextMessage) {
		fmt.Println("Got a standard message: ", m.Body, " with the id ", m.Id)
	})

	// Message responder
	ec.Subscribe("Messaging.Text.Respond", func(subject, reply string, m *protos.TextMessage) {

		// Respond to the requester using the reply address
		fmt.Println("Asked to respond to a message: ", m.Body, " with the id ", m.Id, ".")
		newMessage := protos.TextMessage{Id: m.Id, Body: "Responding!"}
		ec.Publish(reply, &newMessage)
	})

	select {}
}
