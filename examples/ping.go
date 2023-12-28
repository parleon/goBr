package main

import (
	"log"
	"net/http"

	"github.com/parleon/goBr"
)


func main() {
	app := gobr.NewRuntime()

	pingService := app.NewService()

	pingService.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	pingService.BroadcastService(":8080")

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


}