package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	capi "github.com/hashicorp/consul/api"
	"github.com/silviopilato/gbk8s/pkg/proto"
)

type AgentClient struct {
	rpcClient proto.AgentServiceClient
}

func (a AgentClient) StartWorkload(ctx *gin.Context) {
	var reqBody POSTBody
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	consulClient, err :=capi.NewClient(capi.DefaultConfig())
	cAgent := consulClient.Agent()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	wload := proto.Workload{Image: reqBody.Image, Name: reqBody.Name, PortBindings: reqBody.PortBindings}
	res, err := a.rpcClient.StartWorkload(ctx.Request.Context(), &wload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cAgent.ServiceRegister(&capi.AgentServiceRegistration{Name: reqBody.Name, Port: 8888})
	ctx.JSON(http.StatusOK, res)
}

func (a AgentClient) RemoveWorkload(ctx *gin.Context) {
	name := ctx.Param("name")
	wload := proto.Workload{Name: name}
	consulClient, err :=capi.NewClient(capi.DefaultConfig())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cAgent := consulClient.Agent()

	res, err := a.rpcClient.RemoveWorkload(ctx.Request.Context(), &wload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cAgent.ServiceDeregister(name)	
	ctx.JSON(http.StatusOK, res)
}
