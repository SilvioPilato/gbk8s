package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
	"github.com/silviopilato/gbk8s/pkg/tasks"
)
var cli *client.Client
func main() {
	cli = getDockerClient()
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	log.Printf("Agent started")
	nc.QueueSubscribe("start_workload", "agent", handleStartWorkload)
	nc.QueueSubscribe("delete_workload", "agent", handleRemoveWorkload)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	nc.Drain()
	log.Fatalf("Exiting")
}

func containerStart(containerName string) error {
	return getDockerClient().ContainerStart(context.Background(), containerName, types.ContainerStartOptions{})
}

func containerCreate(imageName string, containerName string) error {
	_, err := getDockerClient().ContainerCreate(context.Background(), &container.Config{Image: imageName}, nil, nil, nil, containerName)
	return err
}

func containerStop(containerName string) error {
	return getDockerClient().ContainerStop(context.Background(), containerName, nil)
}

func containerRemove(containerName string) error {
	return getDockerClient().ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{})
}

func pullImage(imageName string) error {
	reader, err := getDockerClient().ImagePull(context.Background(), imageName, types.ImagePullOptions{})

	if err != nil {
		log.Println(err)
	}
	io.Copy(os.Stdout, reader)
	return err
}

func getDockerClient() *client.Client {
	if cli == nil {
		var err error
		cli, err = client.NewClientWithOpts(client.FromEnv);
		if err != nil {
			panic(err)
		}
	}
	return cli
}

func handleStartWorkload(msg *nats.Msg) {
	task := tasks.ReadTaskFromJSON(msg.Data)
	imageName := task.Properties.Image
	containerName := task.Properties.ContainerName

	log.Printf("Pulling image...")
	err := pullImage(imageName)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Creating container...")
	err = containerCreate(imageName, containerName)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Starting container...")
	err = containerStart(containerName)
	if err != nil {
		log.Println(err)
		return;
	}
	log.Printf("Everything set up!")
}

func handleRemoveWorkload(msg *nats.Msg) {
	task := tasks.ReadTaskFromJSON(msg.Data)
	containerName := task.Properties.ContainerName
	log.Printf("Stopping container...")
	err:= containerStop(containerName)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Removing Container...")
	err= containerRemove(containerName)
	if err != nil {
		log.Println(err)
		return
	} 
	log.Printf("Container removed successfully!")
}