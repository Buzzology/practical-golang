package worker

import (
	"bytes"
	"fmt"
	protosFileServer "github.com/Buzzology/practical-golang/fileserver/protos/fileserver"
	protosTask "github.com/Buzzology/practical-golang/masterworker/protos/masterworker"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var nc *nats.Conn

func main() {

	// Ensure that we've received the nats address
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}

	var err error

	// Connect to nats
	nc, err = nats.Connect(os.Args[1])
	if err != nil {
		fmt.Errorf("Failed to connect to nats: %v", os.Args[1])
		return
	}

	for i := 0; i < 8; i++ {
		go doWork()
	}

	// Block forever
	select {}
}

func doWork() {

	for {

		// Request a task with a 1 second timeout
		msg, err := nc.Request("Work.TaskToDo", nil, 1*time.Second)
		if err != nil {
			fmt.Print("Something went wrong. Waiting two seconds before retrying:", err)
			continue
		}

		// Unmarshal the task
		curTask := protosTask.Task{}
		err = proto.Unmarshal(msg.Data, &curTask)
		if err != nil {
			fmt.Errorf("unable to unmarshal task: %v", err)
			continue
		}

		// Retrieve the file server address
		msg, err = nc.Request("Discovery.FileServer", nil, 1*time.Second)
		if err != nil {
			fmt.Errorf("unable to retrieve file server address: %v", err)
			continue
		}

		// Unmarshal the file server address
		var fileServerAddressTransport = protosFileServer.DiscoverableServiceTransport{}
		err = proto.Unmarshal(msg.Data, &fileServerAddressTransport)
		if err != nil {
			fmt.Errorf("unable to unmarshal file server address: %v", err)
			continue
		}

		// Retrieve file from the file server
		fileServerAddress := fileServerAddressTransport.Address
		resp, err := http.Get(fileServerAddress + "/" + curTask.Uuid)
		if err != nil {
			fmt.Errorf("failed to retrieve file: %v", err)
			continue
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("Failed to read request body: %v", err)
			continue
		}

		// Retrieve and count words in the request body
		words := strings.Split(string(data), ",")
		sort.Strings(words)
		wordCounts := make(map[string]int)
		for i := 0; i < len(words); i++ {
			wordCounts[words[i]] = wordCounts[words[i]] + 1
		}

		resultData := make([]byte, 0, 1024)
		buf := bytes.NewBuffer(resultData)

		// Print results to buffer
		for key, value := range wordCounts {
			fmt.Fprintln(buf, key, ":", value)
		}

		// Generate a new uuid for the finished file
		curTask.Finisheduuid = uuid.NewV4().String()
		resp, err = http.Post(fileServerAddress+"/"+curTask.Finisheduuid, "", buf)
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Errorf("unable to post finished status: %v", err)
			continue
		}

		// Marshal the current task
		data, err = proto.Marshal(&curTask)
		if err != nil {
			fmt.Errorf("failed to marshal task: %v", err)
			continue
		}

		// Notify master worker about finishing the task
		nc.Publish("Work.TaskFinished", data)
	}
}
