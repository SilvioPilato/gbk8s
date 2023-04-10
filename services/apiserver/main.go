package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/silviopilato/gbk8s/pkg/proto"
	"google.golang.org/grpc"
)

type POSTBody struct {
	Image string `json:"image" binding:"required"`
	Name  string `json:"name" binding:"required"`
	PortBindings []*proto.Portbinding `json:"portBindings,omitempty"`
}

type PortBinding struct {
	Outer uint32 `json:"outer"`
	Inner uint32 `json:"inner"`
	Protocol proto.PortProtocol `json:"protocol"`
}

func main() {
	opts := grpc.WithInsecure()
	router := gin.Default()
	conn, err := grpc.Dial(":36666", opts)
	if err != nil {
		log.Fatalf("Error Dialing GRPC agent: %v", err)
	}
	if err != nil {
		log.Fatalf("Error Dialing GRPC agent: %v", err)
	}
	agentClient := AgentClient{rpcClient: proto.NewAgentServiceClient(conn)}
	router.POST("/", agentClient.StartWorkload)
	router.DELETE("/:name", agentClient.RemoveWorkload)
	router.Run()
}
