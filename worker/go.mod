module github.com/Buzzology/practical-golang/worker

replace github.com/Buzzology/practical-golang/masterworker => ../masterworker

replace github.com/Buzzology/practical-golang/worker => ./worker

replace github.com/Buzzology/practical-golang/fileserver => ../fileserver

go 1.16

require (
	github.com/Buzzology/practical-golang/masterworker v0.0.0-00010101000000-000000000000
	github.com/nats-io/nats.go v1.11.0
)
