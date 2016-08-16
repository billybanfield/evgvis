package main

import (
	"time"

	"github.com/billybanfield/evgvis/datamanager"
	"github.com/billybanfield/evgvis/server"
)

func main() {
	go func() {
		for {
			datamanager.UpdateState()
			time.Sleep(time.Second * 30)
		}
	}()

	server.RunWebServer()
}
