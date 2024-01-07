package main

import (
	"log"
	"net/http"

	"github.com/parleon/goBr"
	"github.com/parleon/goBr/middleware"
)


func main() {
	app := gobr.NewRuntime()

	pingService := app.NewService()


	pingService.Use(middleware.Recovery)

	pingService.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	pingService.BroadcastService("8080")

	app.ServeAllPorts()

	c := app.NewClient()

	req, _ := http.NewRequest("GET", "http://localhost:8080/ping", nil)

	resp, err := c.Do(req)

	if err != nil {
		log.Println(err)
	}

	log.Println(resp)

	extc := http.DefaultClient

	extResp, err := extc.Do(req)

	if err != nil {
		log.Println(err)
	}

	log.Println(extResp)

	app.CloseAllPorts()

	c2 := app.NewClient()

	req2, _ := http.NewRequest("GET", "http://localhost:8080/ping", nil)

	resp2, err := c2.Do(req2)

	if err != nil {
		log.Println(err)
	}

	log.Println(resp2)

	extc2 := http.DefaultClient

	extResp2, err := extc2.Do(req)

	if err != nil {
		log.Println(err)
	}

	log.Println(extResp2)

}