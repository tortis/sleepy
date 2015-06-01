package main

import (
	"fmt"
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
	//api.Filter(apiLogFilter)
	log.Fatal(http.ListenAndServe(":3000", api))
}

func apiLogFilter(w http.ResponseWriter, r *http.Request, d map[string]interface{}) error {
	fmt.Printf("REQUEST: [client %s] to [%s] %s\n", r.RemoteAddr, r.Method, r.URL)
	return nil
}
