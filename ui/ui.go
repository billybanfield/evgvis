package ui

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/billybanfield/heroku2/datamanager"
	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("%#v\n", datamanager.FetchState()))
	log.Println("Page load complete")
}

func RunWebServer() {
	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/", hello)
	router.Run(":" + port)
}
