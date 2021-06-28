package main

import (
	"bytes"
	"fmt"
	protosFileServer "github.com/Buzzology/practical-golang/fileserver/protos/fileserver"
	protos "github.com/Buzzology/practical-golang/masterworker/protos/masterworker"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
	"sync"
	"time"
)

var Tasks []protos.Task
var TaskMutex sync.Mutex
var oldestFinishedTaskPointer int
var nc *nats.Conn

func main() {

	// Ensure that the nats server address is provided
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}

	var err error

	// Try to connect to nats
	nc, err = nats.Connect(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	// Initialise our slice to keep track of our tasks
	Tasks = make([]protos.Task, 0, 20)
	TaskMutex = sync.Mutex{}
	oldestFinishedTaskPointer = 0

	initTestTasks()

	//wg := sync.WaitGroup{}

	// Monitor the Work.TaskToDo subject for any tasks that we should be handling
	nc.Subscribe("Work.TaskToDo", func(m *nats.Msg) {

		myTaskPointer, ok := getNextTask()
		if ok {
			data, err := proto.Marshal(myTaskPointer)
			if err == nil {
				nc.Publish(m.Reply, data)
			}
		}
	})

	// Monitor the Work.TaskFinished subject for any tasks that have finished
	nc.Subscribe("Work.TaskFinished", func(m *nats.Msg) {

		// Unmarshal the task
		myTask := protos.Task{}
		err := proto.Unmarshal(m.Data, &myTask)
		if err != nil {
			fmt.Errorf("failed to unmarshal task: %v", err)
			return
		}

		// Mark the task in the slice as finished
		TaskMutex.Lock()
		Tasks[myTask.Id].State = 2
		Tasks[myTask.Id].Finisheduuid = myTask.Finisheduuid
		TaskMutex.Unlock()
	})

	select {}
}

func initTestTasks() {

	// Create our 20 tests tasks
	for i := 0; i < 20; i++ {

		newTask := protos.Task{Uuid: uuid.NewV4().String(), State: 0}

		// Retrieve the file server's address
		fileServerAddressTransport := protosFileServer.DiscoverableServiceTransport{}
		msg, err := nc.Request("Discovery.Fileserver", nil, 1000*time.Millisecond)

		// Process the response
		if err == nil && msg != nil {
			err := proto.Unmarshal(msg.Data, &fileServerAddressTransport)
			if err != nil {
				fmt.Errorf("failed to unmarshal server address: %v", err)
				continue
			}
		}

		// Send a payload to the file server
		fileServerAddress := fileServerAddressTransport.Address
		data := make([]byte, 0, 1024)
		buf := bytes.NewBuffer(data)
		fmt.Fprint(buf, "get,my,data,my,get,get,have")
		r, err := http.Post(fileServerAddress+"/"+newTask.Uuid, "", buf)

		// Terminate if we hit any errors
		if err != nil || r.StatusCode != http.StatusOK {
			fmt.Errorf("failed to send payload to server: %v %v", err, r.StatusCode)
			continue
		}

		// Add the task to our slice to track it
		newTask.Id = int32(len(Tasks))
		Tasks = append(Tasks, newTask)
	}
}

func getNextTask() (*protos.Task, bool) {

	TaskMutex.Lock()
	defer TaskMutex.Unlock()

	for i := oldestFinishedTaskPointer; i < len(Tasks); i++ {

		// 2 == Finished
		if i == oldestFinishedTaskPointer && Tasks[i].State == 2 {
			oldestFinishedTaskPointer++
		} else {

			if Tasks[i].State == 0 {
				Tasks[i].State = 1

				go resetTaskIfNotFinished(i)
				return &Tasks[i], true
			}
		}
	}

	return nil, false
}

// Task is reset if not finished after two minutes
func resetTaskIfNotFinished(i int) {
	time.Sleep(2 * time.Minute)
	TaskMutex.Lock()

	if Tasks[i].State != 2 {
		Tasks[i].State = 0
	}
}
