package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/silviopilato/gbk8s/pkg/config"
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
	args := os.Args[1:]
	configPath := args[0]
	cfg := config.Config{}
	err := config.ReadYamlConfig(&cfg, configPath)

	if err != nil {
		log.Fatalf("Error Reading configuration: %v", err)
	}
	
	opts := grpc.WithInsecure()
	router := gin.Default()
	httpAddress := fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port)
	agentAddress := fmt.Sprintf("%s:%d", cfg.Agent.Host, cfg.Agent.Port)

	conn, err := grpc.Dial(agentAddress, opts)
	if err != nil {
		log.Fatalf("Error Dialing GRPC agent: %v", err)
	}
	if err != nil {
		log.Fatalf("Error Dialing GRPC agent: %v", err)
	}
	agentClient := AgentClient{rpcClient: proto.NewAgentServiceClient(conn)}
	router.POST("/", agentClient.StartWorkload)
	router.DELETE("/:name", agentClient.RemoveWorkload)
	router.Run(httpAddress)
}
