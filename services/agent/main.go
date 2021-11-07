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

func main() {
	cli := getClient()
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	log.Printf("Agent started")
	nc.QueueSubscribe("tasks", "agent", func(msg *nats.Msg) {
		task := tasks.ReadTaskFromJSON(msg.Data)
		imageName := task.Properties.Image
		containerName := task.Properties.ContainerName

		log.Printf("Pulling image...")
		pullImage(cli, imageName)

		log.Printf("Creating container...")
		containerCreate(cli, imageName, containerName)

		log.Printf("Starting container...")
		containerStart(cli, containerName)

		log.Printf("Everything set up!")
	})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	nc.Drain()
	log.Fatalf("Exiting")
}

func containerStart(cli *client.Client, containerName string) {
	err := cli.ContainerStart(context.Background(), containerName, types.ContainerStartOptions{})

	if err != nil {
		panic(err)
	}
}

func containerCreate(cli *client.Client, imageName string, containerName string) {
	_, err := cli.ContainerCreate(context.Background(), &container.Config{Image: imageName}, nil, nil, nil, containerName)

	if err != nil {
		panic(err)
	}
}

func pullImage(cli *client.Client, imageName string) {
	reader, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})

	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
}

func getClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv);
	if err != nil {
		panic(err)
	}
	return cli
}