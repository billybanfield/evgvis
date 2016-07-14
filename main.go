package main

import (
	//	"time"
	"fmt"

	//	"github.com/billybanfield/heroku2/datamanager"
	"github.com/billybanfield/heroku2/jsonfetcher"
	//	"github.com/billybanfield/heroku2/server"
)

func main() {
	/*
		go func() {
			for {
				datamanager.UpdateState()
				time.Sleep(time.Second * 30)
			}
		}()
	*/

	fmt.Println(jsonfetcher.FetchHosts())
	//server.RunWebServer()
}
