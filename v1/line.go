package v1

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func LINEClientMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cloud Run の環境変数から取得している. Cloud Runで環境変数の設定方法を調べといた方がいい
		bot, err := linebot.New(
			os.Getenv("CHANNEL_SECRET"),
			os.Getenv("CHANNEL_TOKEN"),
		)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		// clientをハンドラー間で共有するためにコンテキストに登録
		ctx = SetLINEClientCtx(ctx, bot)
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
		ctx = SetLINEEventsCtx(ctx, events)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetLINEClientCtx(ctx context.Context, bot *linebot.Client) context.Context {
	return context.WithValue(ctx, linebot.Client{}, bot)
}

func SetLINEEventsCtx(ctx context.Context, events []*linebot.Event) context.Context {
	return context.WithValue(ctx, "events", events)
}

func GetLINEEvents(ctx context.Context) ([]*linebot.Event, error) {
	events, ok := ctx.Value("events").([]*linebot.Event)
	if !ok {
		return nil, errors.New("events not found")
	}
	return events, nil
}

func GetLINEClient(ctx context.Context) (*linebot.Client, error) {
	bot, ok := ctx.Value(linebot.Client{}).(*linebot.Client)
	if !ok {
		return nil, errors.New("client bot not found")
	}
	return bot, nil
}
