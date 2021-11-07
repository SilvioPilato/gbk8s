package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type POSTBody struct {
	Image string `json:"image" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func main() {
	router := gin.Default()
	router.POST("/", func(ctx *gin.Context) {
		var reqBody POSTBody
		err := ctx.BindJSON(&reqBody)
		
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(200, gin.H{
				"submitted": true,
			})
		}
		
	})
	router.Run()
}