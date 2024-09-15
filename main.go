package main

import (
	"log"
	"net/http"

	"github.com/Agustincou/go-cucumber-example/handlers"
)

func main() {

	http.HandleFunc("/users/1", handlers.GetUserHandler)

	http.HandleFunc("/ping", handlers.GetPingHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
