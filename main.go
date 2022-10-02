package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	http.HandleFunc("/", pingHandler)

	lineHandler := http.HandlerFunc(LINEWebhookHandler)
	http.Handle("/callback", LINEClientMiddleware(lineHandler))

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

func pingHandler(w http.ResponseWriter, r *http.Request) {

	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

func LINEWebhookHandler(w http.ResponseWriter, r *http.Request) {
	events, e := GetLINEEvents(r.Context())
	if e == false {
		log.Print("no event")
		return
	}
	bot, e := GetLINEClient(r.Context())
	if e == false {
		log.Panicln("no client")
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func LINEClientMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bot, err := linebot.New(
			os.Getenv("CHANNEL_SECRET"),
			os.Getenv("CHANNEL_TOKEN"),
		)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// clientをハンドラー間で共有するためにコンテキストに登録
		SetLINEClientCtx(r.Context(), bot)
		events, err := bot.ParseRequest(r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// イベントもハンドラー間で共有するためにコンテキストに登録。※イベントは正直コンテキストに入れず、橋渡ししてもいいかも知らない
		SetLINEEventsCtx(r.Context(), events)

		next.ServeHTTP(w, r)
	})
}

func SetLINEClientCtx(ctx context.Context, bot *linebot.Client) context.Context {
	return context.WithValue(ctx, linebot.Client{}, bot)
}

func SetLINEEventsCtx(ctx context.Context, events []*linebot.Event) context.Context {
	return context.WithValue(ctx, []linebot.Event{}, events)
}

func GetLINEEvents(ctx context.Context) ([]linebot.Event, bool) {
	events, err := ctx.Value([]linebot.Event{}).([]linebot.Event)
	return events, err
}

func GetLINEClient(ctx context.Context) (linebot.Client, bool) {
	bot, err := ctx.Value(linebot.Client{}).(linebot.Client)
	return bot, err
}
