package main

import (
	"log"

	"github.com/gin-gonic/gin"
	workload "github.com/silviopilato/gbk8s/pkg/proto"
	"google.golang.org/grpc"
)

type POSTBody struct {
	Image string `json:"image" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

func main() {
	opts := grpc.WithInsecure()
	router := gin.Default()
	conn, err := grpc.Dial(":36666", opts)
	if err != nil {
		log.Fatalf("Error Dialing GRPC agent: %v", err)
	}
	agentClient := AgentClient{rpcClient: workload.NewAgentServiceClient(conn)}
	router.POST("/", agentClient.StartWorkload)
	router.DELETE("/:name", agentClient.RemoveWorkload)
	router.Run()
}
