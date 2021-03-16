package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"vax/cmd/authority/api"
)

func main() {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "8083"
	}

	var ACTOR_ID string
	if ACTOR_ID = os.Getenv("ACTOR_ID"); PORT == "" {
		panic("missing ACTOR_ID ENV")
	}

	server := api.NewServer(context.Background(), ACTOR_ID, "/api/v1")
	http.Handle("/", server)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Printf("problem running authority server '%s': %+v", ACTOR_ID, err)
	}
}
