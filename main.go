package main

import (
	"log"
	"net/http"
	"os"

	linecon "github.com/lll-lll-lll-lll/lineconnpass/src"
)

func main() {
	lineHandler := http.HandlerFunc(linecon.LINEWebhookHandler)
	http.Handle("/callback", linecon.LINEClientMiddleware(lineHandler))
	Run()
}

func Run() {
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println(err)
		return
	}
}
