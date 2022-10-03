package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	linecon "github.com/lll-lll-lll-lll/lineconnpass/v1"

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
	user := linecon.GetUserName()
	if user == "" {
		log.Println(fmt.Errorf("no client: %s", e))
		return
	}
	conpass.ConnpassUSER = user

	query := map[string]string{"nickname": conpass.ConnpassUSER}
	q := linecon.CreateQuery(query)
	conpass.Query = q
	u, err := conpass.CreateURL(conpass.Query)
	if err != nil {
		log.Println(fmt.Errorf("no client: %s", e))
		return
	}
	res, err := conpass.Request(u)
	if err != nil {
		log.Println(fmt.Errorf("no client: %s", e))
		return
	}
	defer res.Body.Close()
	if err := conpass.SetResponse(res); err != nil {
		log.Println(fmt.Errorf("no client: %s", e))
		return
	}

	titles := conpass.ConnpassResponse.GetGroupTitles()

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(titles[0])).Do(); err != nil {
					log.Println(fmt.Errorf("no reply message: %v, message: %v", err, message))
				}
			}
		}
	}
}
