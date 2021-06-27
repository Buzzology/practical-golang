package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	//"github.com/cube2222/Blog/NATS/FrontendBackend"
	"github.com/golang/protobuf/proto"
	"os"
	//"sync"
	protos "github.com/Buzzology/practical-golang/time-provider/protos/transport"
	"time"
)

// Tutorial mentions that you normally shouldn't use globals for this
var nc *nats.Conn

func main() {

	// Ensure that we've received that nats server address
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}

	// Attempt to connect
	var err error
	nc, err = nats.Connect(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	// Subscribe to the TimeTellers queue and use the replyWithTime handler when a message is received
	nc.QueueSubscribe("TimeTeller", "TimeTellers", replyWithTime)
	select {}
}

// replyWithTime uses the reply queue to return the current time.
func replyWithTime(m *nats.Msg) {

	// Prepare a reply containing the current time
	var myReply = protos.Time{
		Time: time.Now().Format(time.RFC3339),
	}

	// Marshal the reply
	data, err := proto.Marshal(&myReply)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Reply to the request subject (accessed via m.Reply)
	fmt.Println("Replying to ", m.Reply)
	nc.Publish(m.Reply, data)
}
