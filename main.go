package main

import (
	"time"

	"github.com/billybanfield/heroku2/datamanager"
	"github.com/billybanfield/heroku2/server"
)

func main() {
	go func() {
		for {
			datamanager.UpdateState()
			time.Sleep(time.Minute)
		}
	}()

	server.RunWebServer()
}
