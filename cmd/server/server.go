package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sanderhahn/htpasswd/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}
	log.Printf("Started listening on http://localhost:%s/", port)
	http.HandleFunc("/api", api.DispatchHandler)
	// http.Handle("/", http.FileServer(http.Dir("./public")))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
