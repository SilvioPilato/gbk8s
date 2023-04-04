package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	workload "github.com/silviopilato/gbk8s/pkg/proto"
)

type AgentClient struct {
	rpcClient workload.AgentServiceClient
}

func (a AgentClient) StartWorkload(ctx *gin.Context) {
	var reqBody POSTBody
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wload := workload.Workload{Image: reqBody.Image, Name: reqBody.Name}
	res, err := a.rpcClient.StartWorkload(ctx.Request.Context(), &wload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (a AgentClient) RemoveWorkload(ctx *gin.Context) {
	name := ctx.Param("name")
	wload := workload.Workload{Name: name}
	res, err := a.rpcClient.RemoveWorkload(ctx.Request.Context(), &wload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}
