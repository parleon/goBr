package main

import (	
	gobr "github.com/parleon/goBr"
	reg "github.com/parleon/goBr/examples/regression/services"
)

func main() {
	app := gobr.NewRuntime()
	reg.RegressionService(app, "8080")
	app.Open()
	for {
		
	}
}