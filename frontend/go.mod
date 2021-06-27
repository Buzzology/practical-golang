module github.com/Buzzology/practical-golang/frontend

replace github.com/Buzzology/practical-golang/frontend => ./frontend

replace github.com/Buzzology/practical-golang/provider => ../provider

replace github.com/Buzzology/practical-golang/time-provider => ../time-provider

go 1.16

require (
	github.com/Buzzology/practical-golang/provider v0.0.0-00010101000000-000000000000
	github.com/Buzzology/practical-golang/time-provider v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.5.0
	github.com/gorilla/mux v1.8.0
	github.com/nats-io/nats.go v1.11.0
)
