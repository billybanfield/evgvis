package server

import (
	"log"
	"net/http"
	"os"

	"github.com/billybanfield/heroku2/datamanager"
	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
	log.Println("Page load complete")
}

func fetchData(c *gin.Context) {
	c.JSON(http.StatusOK, datamanager.FetchState())
	log.Println("JSON load complete")
}

func RunWebServer() {
	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())

	router.LoadHTMLGlob("server/templates/*.tmpl.html")
	router.Static("/static", "server/static/")

	router.GET("/", index)
	router.GET("/data", fetchData)
	router.Run(":" + port)
}
