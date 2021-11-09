package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/silviopilato/gbk8s/pkg/tasks"
)

type POSTBody struct {
	Image string `json:"image" binding:"required"`
	Name string `json:"name" binding:"required"`
}

var nc *nats.Conn
func main() {
	router := gin.Default()
	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	router.POST("/", PostHandler)
	router.DELETE("/:name", DeleteHandler)
	router.Run()
}

func PostHandler(ctx *gin.Context) {
	var reqBody POSTBody
	err := ctx.BindJSON(&reqBody)
	
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	properties := tasks.TaskProperties{Image: reqBody.Image, ContainerName: reqBody.Name}
	workload := tasks.GetStartWorkloadTask(properties)
	serialized := tasks.SerializeTask(&workload)
	err = nc.Publish("start_workload", serialized)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(200, gin.H{
		"submitted": true,
	})
}

func DeleteHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	properties := tasks.TaskProperties{ContainerName: name}
	workload := tasks.GetRemoveWorkloadTask(properties)
	serialized := tasks.SerializeTask(&workload)
	err := nc.Publish("delete_workload", serialized)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(200, gin.H{
		"submitted": true,
	})
}