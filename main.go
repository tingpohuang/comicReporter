// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				_, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				log.Println(message.Text)

				jsonData := []byte(`{
					"type": "bubble",
					"hero": {
					  "type": "image",
					  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_1_cafe.png",
					  "size": "full",
					  "aspectRatio": "20:13",
					  "aspectMode": "cover",
					  "action": {
						"type": "uri",
						"uri": "http://linecorp.com/"
					  }
					},
					"body": {
					  "type": "box",
					  "layout": "vertical",
					  "contents": [
						{
						  "type": "text",
						  "text": "Brown Cafe",
						  "weight": "bold",
						  "size": "xl"
						},
						{
						  "type": "box",
						  "layout": "baseline",
						  "margin": "md",
						  "contents": [
							{
							  "type": "icon",
							  "size": "sm",
							  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
							},
							{
							  "type": "icon",
							  "size": "sm",
							  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
							},
							{
							  "type": "icon",
							  "size": "sm",
							  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
							},
							{
							  "type": "icon",
							  "size": "sm",
							  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
							},
							{
							  "type": "icon",
							  "size": "sm",
							  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gray_star_28.png"
							},
							{
							  "type": "text",
							  "text": "4.0",
							  "size": "sm",
							  "color": "#999999",
							  "margin": "md",
							  "flex": 0
							}
						  ]
						},
						{
						  "type": "box",
						  "layout": "vertical",
						  "margin": "lg",
						  "spacing": "sm",
						  "contents": [
							{
							  "type": "box",
							  "layout": "baseline",
							  "spacing": "sm",
							  "contents": [
								{
								  "type": "text",
								  "text": "Place",
								  "color": "#aaaaaa",
								  "size": "sm",
								  "flex": 1
								},
								{
								  "type": "text",
								  "text": "Miraina Tower, 4-1-6 Shinjuku, Tokyo",
								  "wrap": true,
								  "color": "#666666",
								  "size": "sm",
								  "flex": 5
								}
							  ]
							},
							{
							  "type": "box",
							  "layout": "baseline",
							  "spacing": "sm",
							  "contents": [
								{
								  "type": "text",
								  "text": "Time",
								  "color": "#aaaaaa",
								  "size": "sm",
								  "flex": 1
								},
								{
								  "type": "text",
								  "text": "10:00 - 23:00",
								  "wrap": true,
								  "color": "#666666",
								  "size": "sm",
								  "flex": 5
								}
							  ]
							}
						  ]
						}
					  ]
					},
					"footer": {
					  "type": "box",
					  "layout": "vertical",
					  "spacing": "sm",
					  "contents": [
						{
						  "type": "button",
						  "style": "link",
						  "height": "sm",
						  "action": {
							"type": "uri",
							"label": "CALL",
							"uri": "https://linecorp.com"
						  }
						},
						{
						  "type": "button",
						  "style": "link",
						  "height": "sm",
						  "action": {
							"type": "uri",
							"label": "WEBSITE",
							"uri": "https://linecorp.com"
						  }
						},
						{
						  "type": "spacer",
						  "size": "sm"
						}
					  ],
					  "flex": 0
					}
				  }`)

				container, err := linebot.UnmarshalFlexMessageJSON(jsonData)
				// err is returned if invalid JSON is given that cannot be unmarshalled
				if err != nil {
					log.Print(err)
				}
				flexMessage := linebot.NewFlexMessage("alt text", container)

				if _, err = bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
					log.Print(err)
				}

			}
		}
	}
}
