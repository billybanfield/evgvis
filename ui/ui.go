package ui

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/billybanfield/heroku2/jsonfetcher"
	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	hosts := jsonfetcher.FetchPage()
	c.String(http.StatusOK, fmt.Sprintf("%#v", hosts))
	log.Println("page load complete")
}

func RunWebServer() {
	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/", hello)
	router.Run(":" + port)
}
