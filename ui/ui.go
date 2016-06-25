package ui

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.String(http.StatusOK, "works")
	log.Println("page load complete")
}

func RunWebServer() {
	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/", hello)
	router.Run(":" + port)
}
