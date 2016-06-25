package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func hello(c *gin.Context) {
	c.String(http.StatusOK, "hello")
}

func main() {
	port := os.Getenv("PORT")
	router := gin.New()
	router.GET("/", hello)
	router.Run(":" + port)
}
