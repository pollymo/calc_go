package main

import (
	"calc_go/internal"
)

func main() {
	app := internal.NewApplication()
	app.RunServer()
}
