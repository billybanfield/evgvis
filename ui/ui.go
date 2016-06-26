package ui

import (
	"log"
	"net/http"
	"os"

	"github.com/billybanfield/heroku2/datamanager"
	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"data": string(datamanager.FetchState())})
	log.Println("Page load complete")
}

func RunWebServer() {
	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())

	router.LoadHTMLGlob("ui/templates/*.tmpl.html")

	router.GET("/", index)
	router.Run(":" + port)
}
