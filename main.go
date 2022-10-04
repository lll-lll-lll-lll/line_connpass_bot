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
	if err := initRequest(conpass, query); err != nil {
		log.Println(err)
		return
	}

	seriesId := conpass.JoinGroupIdsByComma()
	sm := linecon.GetForThreeMonthsEvent()
	qd := make(map[string]string)
	qd["series_id"] = seriesId
	qd["count"] = "100"
	qd["ym"] = sm
	q := linecon.CreateQuery(qd)
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

	flexsMessages := linecon.CreateConnpassEventFlexMessages(conpass.ConnpassResponse.Events)
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err := bot.ReplyMessage(
					event.ReplyToken,
					flexsMessages[:4]...,
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
	q := linecon.CreateQuery(query)
	c.Query = q
	u, err := c.CreateURL(c.Query)
	if err != nil {
		return err
	}
	res, err := c.Request(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := c.SetResponse(res); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
