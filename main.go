package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	linecon "github.com/lll-lll-lll-lll/lineconnpass/v1"

	"strconv"

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
	flexMessages := CreateConnpassEventFlexMessages(conpass.ConnpassResponse)

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err := bot.ReplyMessage(
					event.ReplyToken,
					flexMessages...,
				).Do(); err != nil {
					log.Println(message)
					return
				}
			}
		}
	}
}

func CreateConnpassEventFlexMessages(connpassResponse *linecon.ConnpassResponse) []linebot.SendingMessage {
	var messages []linebot.SendingMessage
	// flexs := []*linebot.FlexMessage{}
	events := connpassResponse.Events
	for _, e := range events {
		joinedNum := strconv.Itoa(e.Accepted)
		contents := &linebot.BubbleContainer{
			Type: linebot.FlexContainerTypeBubble,
			Body: &linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeVertical,
				Contents: []linebot.FlexComponent{
					&linebot.TextComponent{
						Type:   linebot.FlexComponentTypeText,
						Text:   e.Title,
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
						Action: linebot.NewURIAction("イベントページへ", e.EventUrl),
					},
					&linebot.SeparatorComponent{
						Type: linebot.FlexComponentTypeSeparator,
					},
				},
			},
		}
		d := linebot.NewFlexMessage("Flex message alt text", contents)
		messages = append(messages, d)
	}
	return messages
}
