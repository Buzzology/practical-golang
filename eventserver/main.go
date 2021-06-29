package main

import (
	"fmt"
	protos "github.com/Buzzology/practical-golang/eventserver/protos/eventserver"
	"github.com/nats-io/nats.go"
	natsp "github.com/nats-io/nats.go/encoders/protobuf"
	"os"
	"time"
)

func main() {

	// Ensure that we've been passed the nats connection address
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}

	// Connect to nats
	nc, err := nats.Connect(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	// Create an encoded connection
	ec, err := nats.NewEncodedConn(nc, natsp.PROTOBUF_ENCODER)
	defer ec.Close()

	// Publish five messages
	fmt.Println("Publishing five message...")
	for i := 0; i < 5; i++ {
		myMessage := protos.TextMessage{Id: int32(i), Body: "Hello over standard!"}

		// Send our message over a standard connection
		err := ec.Publish("Messaging.Text.Standard", &myMessage)
		if err != nil {
			fmt.Errorf("failed to send standard message %v", err)
		}
	}

	// Request five responses
	fmt.Println("Requesting five messages...")
	for i := 0; i < 5; i++ {

		myMessage := protos.TextMessage{Id: int32(i), Body: "Hello please respond."}

		// Request our message response
		res := protos.TextMessage{}
		err := ec.Request("Message.Text.Respond", &myMessage, &res, 200*time.Millisecond)
		if err != nil {
			fmt.Errorf("failed to retrieve a message response: %v", err)
			continue
		}

		fmt.Println(res.Body, " with id ", res.Id)
	}

	// Send messages via channels
	fmt.Println("Sending five messages via channel...")
	sendChannel := make(chan *protos.TextMessage)
	ec.BindSendChan("Messaging.Text.Channel", sendChannel)

	for i := 10; i < 15; i++ {
		myMessage := protos.TextMessage{Id: int32(i), Body: "Hello over channel"}
		sendChannel <- &myMessage
	}

	select {}
}
