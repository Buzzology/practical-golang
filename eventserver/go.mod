module github.com/Buzzology/practical-golang/eventserver

replace github.com/Buzzology/practical-golang/eventserver => ./eventserver

go 1.16

require (
	github.com/golang/protobuf v1.5.2
	github.com/nats-io/nats.go v1.11.0
	github.com/satori/go.uuid v1.2.0
	google.golang.org/protobuf v1.27.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
