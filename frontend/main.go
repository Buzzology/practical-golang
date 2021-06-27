package main

import (
	"fmt"
	protosTransport "github.com/Buzzology/practical-golang/provider/protos/transport"
	protosTime "github.com/Buzzology/practical-golang/time-provider/protos/transport"
	"github.com/gorilla/mux"
	"net/http"
	//"github.com/cube2222/Blog/NATS/FrontendBackend"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"os"
	"sync"
	"time"
)

var nc *nats.Conn

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

	// Initialise web server using mux
	m := mux.NewRouter()
	m.HandleFunc("/{id}", handleUserWithTime)
	http.ListenAndServe(":3000", m)
}

func handleUserWithTime(w http.ResponseWriter, r *http.Request) {

	// Prepare the request messages
	vars := mux.Vars(r)
	user := protosTransport.User{Id: vars["id"]}
	curTime := protosTime.Time{}

	//
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)

	// Get the user's name and update our user reference if successful
	go func() {

		// Marshal the user proto
		data, err := proto.Marshal(&user)
		if err != nil || len(user.Id) == 0 {
			fmt.Println(err)
			w.WriteHeader(500)
			fmt.Println("Problem with parsing the user id.")
			return
		}

		// Send the user data and wait no more than 100ms for a response
		msg, err := nc.Request("UserNameById", data, 100*time.Millisecond)
		if err == nil && msg != nil {

			// If we can successfully unmarshal the user update our instance
			userWithName := protosTransport.User{}
			err := proto.Unmarshal(msg.Data, &userWithName)
			if err == nil {
				user = userWithName
			}
		}

		// Signal to the wait group that this goroutine has completed
		waitGroup.Done()
	}()

	// Get the current time and update our `curTime` if successful
	go func() {

		// Send the request to the time service. We don't need any data for this one.
		msg, err := nc.Request("TimeTeller", nil, time.Duration(100*time.Millisecond))
		if err == nil && msg != nil {

			// Unmarshal the response
			receivedTime := protosTime.Time{}
			err := proto.Unmarshal(msg.Data, &receivedTime)
			if err == nil && msg != nil {
				curTime = receivedTime
			}
		}

		// Notify the wait group that we're done so the wait counter can be decremented
		waitGroup.Done()
	}()

	// Wait for both of the goroutines to complete
	waitGroup.Wait()

	// Display output to the user
	fmt.Println(w, "Hello ", user.Name, " with id ", user.Id, " the time is ", curTime.Time, ".")
}
