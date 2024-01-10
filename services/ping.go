package services

import (
	"net/http"

	gobr "github.com/parleon/goBr"
)

func PingService(app *gobr.Runtime, port string) {
	pingService := app.NewService("ping service", port)

	pingService.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
}