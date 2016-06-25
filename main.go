package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func hello(c *gin.Context) {
	c.String(http.StatusOK, "hello")
}

func main() {
	router := gin.New()
	router.GET("/", hello)
	router.Run(":8000")
}
