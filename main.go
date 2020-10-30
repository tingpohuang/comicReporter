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

				// Web Crawler

				//Preprocessd Flex message Json data

				jsonData := []byte(`{
					"type": "carousel",
					"contents": [
					  {
						"type": "bubble",
						"hero": {
						  "type": "image",
						  "size": "full",
						  "aspectRatio": "20:13",
						  "aspectMode": "cover",
						  "url": "https://mhfm1us.cdnmanhua.net/42/41957/20190701151032_180x240_19.jpg"
						},
						"body": {
						  "type": "box",
						  "layout": "vertical",
						  "spacing": "sm",
						  "contents": [
							{
							  "type": "text",
							  "text": message.Text,
							  "wrap": true,
							  "weight": "bold",
							  "size": "xl"
							},
							{
							  "type": "box",
							  "layout": "baseline",
							  "contents": [
								{
								  "type": "text",
								  "text": "$49",
								  "wrap": true,
								  "weight": "bold",
								  "size": "xl",
								  "flex": 0
								},
								{
								  "type": "text",
								  "text": ".99",
								  "wrap": true,
								  "weight": "bold",
								  "size": "sm",
								  "flex": 0
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
							  "style": "primary",
							  "action": {
								"type": "uri",
								"label": "Add to Cart",
								"uri": "https://linecorp.com"
							  }
							},
							{
							  "type": "button",
							  "action": {
								"type": "uri",
								"label": "Add to wishlist",
								"uri": "https://linecorp.com"
							  }
							}
						  ]
						}
					  },
					  {
						"type": "bubble",
						"hero": {
						  "type": "image",
						  "size": "full",
						  "aspectRatio": "20:13",
						  "aspectMode": "cover",
						  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_6_carousel.png"
						},
						"body": {
						  "type": "box",
						  "layout": "vertical",
						  "spacing": "sm",
						  "contents": [
							{
							  "type": "text",
							  "text": "Metal Desk Lamp",
							  "wrap": true,
							  "weight": "bold",
							  "size": "xl"
							},
							{
							  "type": "box",
							  "layout": "baseline",
							  "flex": 1,
							  "contents": [
								{
								  "type": "text",
								  "text": "$11",
								  "wrap": true,
								  "weight": "bold",
								  "size": "xl",
								  "flex": 0
								},
								{
								  "type": "text",
								  "text": ".99",
								  "wrap": true,
								  "weight": "bold",
								  "size": "sm",
								  "flex": 0
								}
							  ]
							},
							{
							  "type": "text",
							  "text": "Temporarily out of stock",
							  "wrap": true,
							  "size": "xxs",
							  "margin": "md",
							  "color": "#ff5551",
							  "flex": 0
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
							  "flex": 2,
							  "style": "primary",
							  "color": "#aaaaaa",
							  "action": {
								"type": "uri",
								"label": "Add to Cart",
								"uri": "https://linecorp.com"
							  }
							},
							{
							  "type": "button",
							  "action": {
								"type": "uri",
								"label": "Add to wish list",
								"uri": "https://linecorp.com"
							  }
							}
						  ]
						}
					  },
					  {
						"type": "bubble",
						"body": {
						  "type": "box",
						  "layout": "vertical",
						  "spacing": "sm",
						  "contents": [
							{
							  "type": "button",
							  "flex": 1,
							  "gravity": "center",
							  "action": {
								"type": "uri",
								"label": "See more",
								"uri": "https://linecorp.com"
							  }
							}
						  ]
						}
					  }
					]
				  }`)

				container, err := linebot.UnmarshalFlexMessageJSON(jsonData)
				// err is returned if invalid JSON is given that cannot be unmarshalled
				if err != nil {
					log.Print(err)
				}

				flexMessage := linebot.NewFlexMessage("alt text", container)

				// Reply Message

				if _, err = bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
					log.Print(err)
				}

			}
		}
	}
}
