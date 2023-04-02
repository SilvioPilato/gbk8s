package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/docker/docker/client"
	workload "github.com/silviopilato/gbk8s/pkg/proto"
	"google.golang.org/grpc"
)

var cli *client.Client

func main() {
	cli = getDockerClient()
	lis, err := net.Listen("tcp", ":36666")
	if err != nil {
		log.Fatalf("Error building tcp listener: %v", err)
	}
	server := grpc.NewServer()
	workload.RegisterAgentServiceServer(server, AgentService{dockerClient: cli})
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Error serving grpc: %v", err)
	}
	log.Printf("Agent started")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Fatalf("Exiting")
}

func getDockerClient() *client.Client {
	if cli == nil {
		var err error
		cli, err = client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			log.Fatalf("Error binding to docker client: %v", err)
		}
	}
	return cli
}
