package main

import (
	"log"

	"github.com/Ranik23/url-shortener/internal/app"
)





func main() {

	app, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to create the app: %v", err)
	}
	app.Run()
}