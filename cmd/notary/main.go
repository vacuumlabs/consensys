package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"vax/cmd/notary/api"
)

func main() {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "8080"
	}

	server := api.NewServer(context.Background(), "/api/v1")
	http.Handle("/", server)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Printf("problem running notary server: %+v", err)
	}
}
