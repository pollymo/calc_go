package main

import (
	"calc_go/internal/application"
)

func main() {
	app := application.NewApplication()
	app.RunServer()
}
