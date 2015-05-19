package main

import (
	"log"
	"net/http"

	"github.com/tortis/sleepy"
)

func main() {
	// Create a new API
	api := sleepy.New("/v2")

	// Init user resource
	userRes := UserResource{}
	api.Register(userRes.Generate())
	log.Fatal(http.ListenAndServe(":8080", api))
}
