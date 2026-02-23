package main

import (
	"net/http"
	"log"

	"jobfinder/internal/app"
)


func main() {
	application, handler := app.New()

	log.Println("server is running on port:", application.Config.Port)

	log.Fatal(http.ListenAndServe(":"+application.Config.Port, handler))
}