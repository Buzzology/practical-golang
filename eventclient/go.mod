module github.com/Buzzology/practical-golang/eventclient

replace github.com/Buzzology/practical-golang/eventclient => ./eventclient

replace github.com/Buzzology/practical-golang/eventserver => ../eventserver

go 1.16

require (
	github.com/Buzzology/practical-golang/eventserver v0.0.0-00010101000000-000000000000
	github.com/nats-io/nats-server/v2 v2.3.0 // indirect
	github.com/nats-io/nats.go v1.11.0
)
