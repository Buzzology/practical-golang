# Practical Golang with Jacob Martin
Using the following tutorials to learn how to implement NATS with gRPC and Golang.

## Tutorials
Article list 
https://jacobmartins.com/article-list/  

Article #1: Practical Golang: Getting started with NATS and related patterns
- Setup NATS


## NATS
Setup using a docker image:
- Retrieve the docker image: `docker pull nats`  
- Run the container: `docker container run nats`

Retrieve the go library: `go get https://github.com/nats-io/nats` 

### QueueSubscribe
A message will be sent to a single subscriber. If there are multiple subscribers on one will receive/process the message.