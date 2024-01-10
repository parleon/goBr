package test

import (
	"net/http"
	"testing"

	"github.com/parleon/goBr"
)

func TestLocalPingService(t *testing.T) {
	app := gobr.NewRuntime()

	pingService := app.NewService("ping service", "8080")

	pingService.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	app.Open()

	c := app.NewClient()

	req, _ := http.NewRequest("GET", "http://localhost:8080/ping", nil)

	resp, err := c.Do(req)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Error("Status code is not 200")
	}

	app.Close()

}

func TestRemotePingService(t *testing.T) {
	app := gobr.NewRuntime()

	pingService := app.NewService("ping service", "8081")

	pingService.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	app.Open()

	req, _ := http.NewRequest("GET", "http://localhost:8081/ping", nil)

	extc := http.DefaultClient

	extResp, err := extc.Do(req)

	if err != nil {
		t.Error(err)
	}

	if extResp.StatusCode != 200 {
		t.Error("Status code is not 200")
	}

	app.Close()
	
}

func BenchmarkLocalPingService(b *testing.B) {

	app := gobr.NewRuntime()

	pingService := app.NewService("ping service","8082")

	pingService.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	app.Open()

	c := app.NewClient()

	req, _ := http.NewRequest("GET", "http://localhost:8082/ping", nil)


	for i := 0; i<b.N; i++ {
		resp, err := c.Do(req)

		if err != nil {
			b.Error(err)
		}

		if resp.StatusCode != 200 {
			b.Error("Status code is not 200")
		}
	}

	app.Close()

}

func BenchmarkRemotePingService(b *testing.B) {
	app := gobr.NewRuntime()

	pingService := app.NewService("ping service", "8083")

	pingService.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	app.Open()

	req, _ := http.NewRequest("GET", "http://localhost:8083/ping", nil)

	extc := http.DefaultClient

	for i := 0; i<b.N; i++ {

	extResp, err := extc.Do(req)

		if err != nil {
			b.Error(err)
		}

		if extResp.StatusCode != 200 {
			b.Error("Status code is not 200")
		}
	}

	app.Close()
	
}

func BenchmarkRemoteHTTPServer(b *testing.B) {

	smux := http.NewServeMux()
	
	smux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	server := http.Server{
		Addr: ":8084",
		Handler: smux,
	}
	go server.ListenAndServe()

	req, _ := http.NewRequest("GET", "http://localhost:8084/ping", nil)

	extc := http.DefaultClient

	for i := 0; i < b.N; i++ {
		extResp, err := extc.Do(req)

		if err != nil {
			b.Error(err)
		}

		if extResp.StatusCode != 200 {
			b.Error("Status code is not 200")
		}
	}

	server.Close()

}