package main

import (
	"fmt"
	protos "github.com/Buzzology/practical-golang/fileserver/protos/fileserver"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"io"
	"net/http"
	"os"
)

// main This service acts as an api for uploading files
func main() {

	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}

	m := mux.NewRouter()

	m.HandleFunc("/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		file, err := os.Open("/tmp/" + vars["name"])
		defer file.Close()
		if err != nil {
			w.WriteHeader(404)
		}
		if file != nil {
			_, err := io.Copy(w, file)
			if err != nil {
				w.WriteHeader(500)
			}
		}
	}).Methods("GET")

	m.HandleFunc("/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		file, err := os.Create("/tmp/" + vars["name"])
		defer file.Close()
		if err != nil {
			w.WriteHeader(500)
		}
		if file != nil {
			_, err := io.Copy(file, r.Body)
			if err != nil {
				w.WriteHeader(500)
			}
		}
	}).Methods("POST")

	RunServiceDiscoverable()

	http.ListenAndServe(":3000", m)
}

func RunServiceDiscoverable() {

	// Ensure that the nats address was provided
	nc, err := nats.Connect(os.Args[1])
	if err != nil {
		fmt.Println("Can't connect to NATS. Service is not discoverable.")
	}

	// Subscribe to the discovery subject and reply to any messages with our address so that they know how to reach us.
	nc.Subscribe("Discovery.FileServer", func(m *nats.Msg) {
		serviceAddressTransport := protos.DiscoverableServiceTransport{Address: "http://localhost:3000"}
		data, err := proto.Marshal(&serviceAddressTransport)
		if err == nil {

			// Reply to the sender
			nc.Publish(m.Reply, data)
		}
	})
}
