module github.com/Buzzology/practical-golang/fileserver

replace github.com/Buzzology/practical-golang/fileserver => ./fileserver

go 1.16

require (
	github.com/golang/protobuf v1.5.0
	github.com/gorilla/mux v1.8.0
	github.com/nats-io/nats-server/v2 v2.3.0 // indirect
	github.com/nats-io/nats.go v1.11.0
	google.golang.org/protobuf v1.27.0 // indirect
)
