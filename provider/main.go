package main

import (
	"fmt"
	protos "github.com/Buzzology/practical-golang/provider/protos/transport"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"os"
)

var nc *nats.Conn
var users map[string]string

func main() {

	// Validate that we've received an address for nats
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address")
		return
	}

	// Connect to nats
	var err error
	nc, err = nats.Connect(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	// Initialise the map for our usernames
	users = make(map[string]string)
	users["1"] = "Batman"
	users["2"] = "Spiderman"
	users["3"] = "Superman"
	users["4"] = "Flash"

	// This will subscribe us to a queue. If multiple instances are subscribed to UserNameById
	// each message will only be received by one instance. This is also async so we will need
	// to create a blocking call
	nc.QueueSubscribe("UserNameById", "userNameByIdProviders", replyWithUserId)

	// Because the queuesubscribe is async we need to block. Select normally used to block until a message is received
	// on a channel. https://github.com/Buzzology/go-freecodecamp#select-statements
	select {}
}

// replyWithUserId retrieves a user from the message and sets the relevant name.
func replyWithUserId(m *nats.Msg) {

	// Unmarshall the user from the message
	myUser := protos.User{}
	err := proto.Unmarshal(m.Data, &myUser)
	if err != nil {
		fmt.Println(err)
	}

	// Retrieve the name from our local map and marshal it
	myUser.Name = users[myUser.Id]
	data, err := proto.Marshal(&myUser)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Reply to the request topic (accessed via m.Reply)
	fmt.Println("Replying to ", m.Reply)
	nc.Publish(m.Reply, data)
}
