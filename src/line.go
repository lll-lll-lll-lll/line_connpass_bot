package v1

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

// CreateConnpassEventFlexMessages cntの数値分だけeventからflex messageを作成する再帰関数
func CreateConnpassEventFlexMessages(e Event) linebot.SendingMessage {
	joinedNum := strconv.Itoa(e.Accepted)
	title := e.Title
	eventURL := e.EventUrl
	message := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   title,
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Align:  "center",
				},
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Margin:  linebot.FlexComponentMarginTypeLg,
					Spacing: linebot.FlexComponentSpacingTypeSm,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "参加者",
									Color: "#aaaaaa",
									Size:  linebot.FlexTextSizeTypeSm,
									Flex:  linebot.IntPtr(1),
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  joinedNum,
									Wrap:  true,
									Color: "#666666",
									Size:  linebot.FlexTextSizeTypeSm,
									Flex:  linebot.IntPtr(3),
								},
							},
						},
					},
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeSm,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypeLink,
					Height: linebot.FlexButtonHeightTypeSm,
					Action: linebot.NewURIAction("イベントページへ", eventURL),
				},
				&linebot.SeparatorComponent{
					Type: linebot.FlexComponentTypeSeparator,
				},
			},
		},
	}
	flexMess := linebot.NewFlexMessage("Flex message alt text", message)
	return flexMess
}

func LINEWebhookHandler(w http.ResponseWriter, r *http.Request) {
	events, e := GetLINEEvents(r.Context())
	if e != nil {
		log.Println(fmt.Errorf("no events: %s", e))
		return
	}
	bot, e := GetLINEClient(r.Context())
	if e != nil {
		log.Println(fmt.Errorf("no client: %s", e))
		return
	}
	conpass := NewConnpass()
	query := map[string]string{"keyword": "go"}
	if err := conpass.Request(conpass, query); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				for _, e := range conpass.ConnpassResponse.Events {
					flexMessage := CreateConnpassEventFlexMessages(e)
					if _, err := bot.ReplyMessage(
						event.ReplyToken,
						flexMessage,
					).Do(); err != nil {
						log.Println(message)
						return
					}
				}

			}
		}
	}
}
