package services

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/parleon/goBr"
	"github.com/parleon/goBr/middleware"
)

func RegressionService(app *gobr.Runtime, port string) {
	regressionService := app.NewService("regression service", port)

	regressionService.Use(middleware.Timer)

	regressionService.HandleFunc("/regress", func(w http.ResponseWriter, r *http.Request) {
		pScriptPath := "./scripts"

		if err := os.Setenv("PYTHONPATH", pScriptPath + "/venv/lib/python3.10/site-packages"); err != nil {
			log.Println(err)
		}
		
		read, write := io.Pipe()

		params, err := url.ParseQuery(r.URL.RawQuery)

		if err != nil {
			log.Println(err)
		}

		if params["y"] == nil{ 
			return
		}



		cmd := exec.Command("python3","./scripts/regress.py", params["y"][0])

		if params["x"] != nil {
			cmd.Args = append(cmd.Args, params["x"]...)
		}

		cmd.Stdout = write
		cmd.Stderr = os.Stderr

		go func() {
			cmd.Run()
			write.Close()
		}()

		io.Copy(w, read)
		
	})
}