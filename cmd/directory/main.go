package main

import (
	"context"
	"net/http"
	"os"
	"vax/cmd/directory/api"
)

func main() {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "8081"
	}

	server := api.NewServer(context.Background(), "/api/v1")
	http.Handle("/", server)
	http.ListenAndServe(":"+PORT, nil)
}
