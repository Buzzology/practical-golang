# Practical Golang with Jacob Martin
Using the following tutorials to learn how to implement NATS with gRPC and Golang.

## Tutorials
Article list 
https://jacobmartins.com/article-list/  

Article #1a: Practical Golang: Getting started with NATS and related patterns
- Setup NATS
- Setup frontend service
- Setup provider service
- Setup time-provider service

Article #1b: The master-slave pattern  
- Setup a file server with service discovery
- Setup a master worker to distribute 20 test tasks and any additional tasks received via NATS
- Setup a worker to process tasks


## NATS
Setup using a docker image:
- Retrieve the docker image: `docker pull nats`  
- Run the container: `docker run -d --name nats-main -p 4222:4222 -p 6222:6222 -p 8222:8222 nats`

**Ports**  
4222: Clients to connect  
6222: Routing port for clustering  
8222: Http overview and management  

Retrieve the go library: `go get https://github.com/nats-io/nats` 

### QueueSubscribe
A message will be sent to a single subscriber. If there are multiple subscribers on one will receive/process the message.

## Running a service
`go run main.go localhost:4222`