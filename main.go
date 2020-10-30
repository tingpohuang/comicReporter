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
					"type": "carousel",
					"contents": [
					  {
						"type": "bubble",
						"body": {
						  "type": "box",
						  "layout": "vertical",
						  "contents": [
							{
							  "type": "image",
							  "url": "https://mhfm1us.cdnmanhua.net/42/41957/20190701151032_180x240_19.jpg",
							  "size": "full",
							  "aspectMode": "cover",
							  "aspectRatio": "2:3",
							  "gravity": "top"
							},
							{
							  "type": "box",
							  "layout": "vertical",
							  "contents": [
								{
								  "type": "box",
								  "layout": "vertical",
								  "contents": [
									{
									  "type": "text",
									  "text": "Brown's T-shirts",
									  "size": "xl",
									  "color": "#ffffff",
									  "weight": "bold"
									}
								  ]
								},
								{
								  "type": "box",
								  "layout": "baseline",
								  "contents": [
									{
									  "type": "text",
									  "text": "¥35,800",
									  "color": "#ebebeb",
									  "size": "sm",
									  "flex": 0
									},
									{
									  "type": "text",
									  "text": "¥75,000",
									  "color": "#ffffffcc",
									  "decoration": "line-through",
									  "gravity": "bottom",
									  "flex": 0,
									  "size": "sm"
									}
								  ],
								  "spacing": "lg"
								},
								{
								  "type": "box",
								  "layout": "vertical",
								  "contents": [
									{
									  "type": "filler"
									},
									{
									  "type": "box",
									  "layout": "baseline",
									  "contents": [
										{
										  "type": "filler"
										},
										{
										  "type": "icon",
										  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/flexsnapshot/clip/clip14.png"
										},
										{
										  "type": "text",
										  "text": "Add to cart",
										  "color": "#ffffff",
										  "flex": 0,
										  "offsetTop": "-2px",
										  "action": {
											"type": "uri",
											"label": "action",
											"uri": "http://linecorp.com/"
										  },
										  "contents": [
											{
											  "type": "span",
											  "text": "hello, world"
											}
										  ]
										},
										{
										  "type": "filler"
										}
									  ],
									  "spacing": "sm"
									},
									{
									  "type": "filler"
									}
								  ],
								  "borderWidth": "1px",
								  "cornerRadius": "4px",
								  "spacing": "sm",
								  "borderColor": "#ffffff",
								  "margin": "xxl",
								  "height": "40px"
								}
							  ],
							  "position": "absolute",
							  "offsetBottom": "0px",
							  "offsetStart": "0px",
							  "offsetEnd": "0px",
							  "backgroundColor": "#03303Acc",
							  "paddingAll": "20px",
							  "paddingTop": "18px"
							},
							{
							  "type": "box",
							  "layout": "vertical",
							  "contents": [
								{
								  "type": "text",
								  "text": "SALE",
								  "color": "#ffffff",
								  "align": "center",
								  "size": "xs",
								  "offsetTop": "3px"
								}
							  ],
							  "position": "absolute",
							  "cornerRadius": "20px",
							  "offsetTop": "18px",
							  "backgroundColor": "#ff334b",
							  "offsetStart": "18px",
							  "height": "25px",
							  "width": "53px"
							}
						  ],
						  "paddingAll": "0px"
						}
					  },
					  {
						"type": "bubble",
						"body": {
						  "type": "box",
						  "layout": "vertical",
						  "contents": [
							{
							  "type": "image",
							  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/flexsnapshot/clip/clip2.jpg",
							  "size": "full",
							  "aspectMode": "cover",
							  "aspectRatio": "2:3",
							  "gravity": "top"
							},
							{
							  "type": "box",
							  "layout": "vertical",
							  "contents": [
								{
								  "type": "box",
								  "layout": "vertical",
								  "contents": [
									{
									  "type": "text",
									  "text": "Cony's T-shirts",
									  "size": "xl",
									  "color": "#ffffff",
									  "weight": "bold"
									}
								  ]
								},
								{
								  "type": "box",
								  "layout": "baseline",
								  "contents": [
									{
									  "type": "text",
									  "text": "¥35,800",
									  "color": "#ebebeb",
									  "size": "sm",
									  "flex": 0
									},
									{
									  "type": "text",
									  "text": "¥75,000",
									  "color": "#ffffffcc",
									  "decoration": "line-through",
									  "gravity": "bottom",
									  "flex": 0,
									  "size": "sm"
									}
								  ],
								  "spacing": "lg"
								},
								{
								  "type": "box",
								  "layout": "vertical",
								  "contents": [
									{
									  "type": "filler"
									},
									{
									  "type": "box",
									  "layout": "baseline",
									  "contents": [
										{
										  "type": "filler"
										},
										{
										  "type": "icon",
										  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/flexsnapshot/clip/clip14.png"
										},
										{
										  "type": "text",
										  "text": "Add to cart",
										  "color": "#ffffff",
										  "flex": 0,
										  "offsetTop": "-2px"
										},
										{
										  "type": "filler"
										}
									  ],
									  "spacing": "sm"
									},
									{
									  "type": "filler"
									}
								  ],
								  "borderWidth": "1px",
								  "cornerRadius": "4px",
								  "spacing": "sm",
								  "borderColor": "#ffffff",
								  "margin": "xxl",
								  "height": "40px"
								}
							  ],
							  "position": "absolute",
							  "offsetBottom": "0px",
							  "offsetStart": "0px",
							  "offsetEnd": "0px",
							  "backgroundColor": "#9C8E7Ecc",
							  "paddingAll": "20px",
							  "paddingTop": "18px"
							},
							{
							  "type": "box",
							  "layout": "vertical",
							  "contents": [
								{
								  "type": "text",
								  "text": "SALE",
								  "color": "#ffffff",
								  "align": "center",
								  "size": "xs",
								  "offsetTop": "3px"
								}
							  ],
							  "position": "absolute",
							  "cornerRadius": "20px",
							  "offsetTop": "18px",
							  "backgroundColor": "#ff334b",
							  "offsetStart": "18px",
							  "height": "25px",
							  "width": "53px"
							}
						  ],
						  "paddingAll": "0px"
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

				if _, err = bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
					log.Print(err)
				}

			}
		}
	}
}
