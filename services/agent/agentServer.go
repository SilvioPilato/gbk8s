package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	workload "github.com/silviopilato/gbk8s/pkg/proto"
)

type AgentService struct {
	*workload.UnimplementedAgentServiceServer
	dockerClient *client.Client
}

func (a AgentService) StartWorkload(ctx context.Context, wload *workload.Workload) (*workload.WorkloadResponse, error) {
	imageName := wload.Image
	containerName := wload.Name

	log.Printf("Pulling image...")
	err := a.pullImage(imageName)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Creating container...")
	err = a.containerCreate(imageName, containerName)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Starting container...")
	err = a.containerStart(containerName)
	if err != nil {
		log.Println(err)
		return &workload.WorkloadResponse{Status: "Failed"}, err
	}
	log.Printf("Everything set up!")
	return &workload.WorkloadResponse{Status: "Success"}, nil
}

func (a AgentService) RemoveWorkload(ctx context.Context, wload *workload.Workload) (*workload.WorkloadResponse, error) {
	containerName := wload.Name
	log.Printf("Stopping container...")
	err := a.containerStop(containerName)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Removing Container...")
	err = a.containerRemove(containerName)
	if err != nil {
		log.Println(err)
		return &workload.WorkloadResponse{Status: "Failed"}, err
	}
	log.Printf("Container removed successfully!")
	return &workload.WorkloadResponse{Status: "Success"}, nil
}

func (a AgentService) containerStart(containerName string) error {
	return a.dockerClient.ContainerStart(context.Background(), containerName, types.ContainerStartOptions{})
}

func (a AgentService) containerCreate(imageName string, containerName string) error {
	_, err := a.dockerClient.ContainerCreate(context.Background(), &container.Config{Image: imageName}, nil, nil, nil, containerName)
	return err
}

func (a AgentService) containerStop(containerName string) error {
	return a.dockerClient.ContainerStop(context.Background(), containerName, nil)
}

func (a AgentService) containerRemove(containerName string) error {
	return a.dockerClient.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{})
}

func (a AgentService) pullImage(imageName string) error {
	reader, err := a.dockerClient.ImagePull(context.Background(), imageName, types.ImagePullOptions{})

	if err != nil {
		log.Println(err)
	}
	io.Copy(os.Stdout, reader)
	return err
}
