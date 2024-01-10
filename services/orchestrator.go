package services

import (
	"encoding/json"
	"log"
	"net/http"

	gobr "github.com/parleon/goBr"
)

type ServiceToggleBody struct {
	ServiceName string
}

func AdminService(app *gobr.Runtime, port string) {
	adminService := app.NewService("admin service", port)

	adminService.HandleFunc("/services/enable", func(w http.ResponseWriter, r *http.Request) {

		var st ServiceToggleBody
		json.NewDecoder(r.Body).Decode(&st)

		app.Services[st.ServiceName].Enabled = true

		log.Println("enabling service: " + st.ServiceName)
	})

	adminService.HandleFunc("/services/disable", func(w http.ResponseWriter, r *http.Request) {
		var st ServiceToggleBody
		json.NewDecoder(r.Body).Decode(&st)

		app.Services[st.ServiceName].Enabled = false

		log.Println("disabling service: " + st.ServiceName)
	})

	adminService.HandleFunc("/runtime/reboot", func(w http.ResponseWriter, r *http.Request) {
		log.Println("rebooting runtime")
		app.Reboot()
	})

}
