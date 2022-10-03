package main

import (
	"log"
	"net/http"
	"os"

	lineV2 "github.com/lll-lll-lll-lll/lineconnpass/v2/line"
)

func main() {
	lineHandler := http.HandlerFunc(lineV2.LINEWebhookHandler)
	http.Handle("/callback", lineV2.LINEClientMiddleware(lineHandler))

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
