package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	linecon "github.com/lll-lll-lll-lll/lineconnpass/pkg"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	lineHandler := http.HandlerFunc(LINEWebhookHandler)
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

func LINEWebhookHandler(w http.ResponseWriter, r *http.Request) {
	events, e := linecon.GetLINEEvents(r.Context())
	if e != nil {
		log.Println(fmt.Errorf("no events: %s", e))
		return
	}
	bot, e := linecon.GetLINEClient(r.Context())
	if e != nil {
		log.Println(fmt.Errorf("no client: %s", e))
		return
	}
	conpass := linecon.NewConnpass()
	query := map[string]string{"keyword": "go"}
	if err := conpass.Request(conpass, query); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	flexsMessages := linecon.CreateConnpassEventFlexMessages(conpass.ConnpassResponse.Events)
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err := bot.ReplyMessage(
					event.ReplyToken,
					flexsMessages[:5]...,
				).Do(); err != nil {
					log.Println(message)
					return
				}
			}
		}
	}
}

// dbを用意してないので、最初にグループIDを取得するためのメソッド
func initRequest(c *linecon.Connpass, query map[string]string) error {

	err := c.Request(c, query)
	if err != nil {
		return err
	}
	return nil
}
