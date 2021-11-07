package main

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/silviopilato/gbk8s/pkg/tasks"
)

func main() {
	cli := getClient()
	jsonFile, err := os.Open("./pkg/tasks/tests/mocks/redis.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	task := tasks.ReadTaskFromJSON(byteValue)

	imageName := task.Properties.Image
	containerName := task.Properties.ContainerName
	pullImage(cli, imageName)
	containerCreate(cli, imageName, containerName)
	containerStart(cli, containerName)
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
	cli, err := client.NewEnvClient();
	if err != nil {
		panic(err)
	}
	return cli
}