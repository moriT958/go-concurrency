package main

import (
	"concurrency/internal/controller"
	"net/http"
)

func main() {
	app := new(controller.AppServer)
	app.Addr = "0.0.0.0:8080"

	mux := http.NewServeMux()
	mux.HandleFunc("POST /convert", app.ConvertGrayscalePost)

	app.Handler = mux

	controller.Start(&app.Server)
}
