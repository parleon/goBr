package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"syscall"
	"testing"

	"github.com/parleon/goBr"
	"github.com/parleon/goBr/services"
)

func TestNameLater(t *testing.T) {
	app := gobr.NewRuntime()

	services.AdminService(app, "9999")

	services.PingService(app, "8080")

	app.Open()

	c := app.NewClient()

	ping_request, _ := http.NewRequest("GET", "http://localhost:8080/ping", nil)

	resp, _ := c.Do(ping_request)

	resp_buf := new(bytes.Buffer)
	resp_buf.ReadFrom(resp.Body)
	log.Println(resp_buf.String())

	s := services.ServiceToggleBody{ServiceName: "ping service"}

	var req_buf bytes.Buffer

	json.NewEncoder(&req_buf).Encode(s)

	req2, _ := http.NewRequest("GET", "http://localhost:9999/services/disable", &req_buf)

	resp2, err := c.Do(req2)

	if err != nil {
		t.Error(err)
	}

	if resp2.StatusCode != 200 {
		t.Error("Status code is not 200")
	}


	reboot_request, _ := http.NewRequest("GET", "http://localhost:9999/runtime/reboot", nil)

	resp3, err := c.Do(reboot_request)

	if err != nil {
		t.Error(err)
	}

	if resp3.StatusCode != 200 {
		t.Error("Status code is not 200")
	}


	_, err = c.Do(ping_request)

	if !errors.Is(err, syscall.ECONNREFUSED) {
		log.Println("service port is still running")
	}

	var buf2 bytes.Buffer

	json.NewEncoder(&buf2).Encode(s)

	req5, _ := http.NewRequest("GET", "http://localhost:9999/services/enable", &buf2)

	resp5, err := c.Do(req5)

	if err != nil {
		t.Error(err)
	}

	if resp5.StatusCode != 200 {
		t.Error("Status code is not 200")
	}

	resp6, err := c.Do(reboot_request)

	if err != nil {
		t.Error(err)
	}

	if resp6.StatusCode != 200 {
		t.Error("Status code is not 200")
	}

	resp7, err := c.Do(ping_request)

	if err != nil {
		t.Error(err)
	}

	if resp7.StatusCode != 200 {
		t.Error("Status code is not 200")
	}

	resp2_buf := new(bytes.Buffer)
	resp2_buf.ReadFrom(resp7.Body)
	log.Println(resp2_buf.String())


}