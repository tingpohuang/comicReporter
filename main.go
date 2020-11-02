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
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

const baseComicURL = "https://www.manhuaren.com"

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
				encodeMessage, _ := url.Parse(message.Text)
				doc, err := goquery.NewDocument("https://www.manhuaren.com/search?title=" + encodeMessage.String())
				if err != nil {
					log.Fatal(err)
				}

				var bubbleContainers []*linebot.BubbleContainer
				doc.Find("ul.book-list li").EachWithBreak(func(index int, item *goquery.Selection) bool {
					book := item
					bookLink, _ := book.Find("a").Attr("href")
					bookTitle, _ := book.Find("a").Attr("title")
					bookImg, _ := book.Find("a img").Attr("src")
					bookInfo := book.Find("p.book-list-info-desc").Contents().Text()
					bubbleContainers = append(bubbleContainers, newBubbleContainer(bookTitle, bookLink, bookImg, bookInfo))

					log.Printf("Link #%d: %s'\n", index, bookLink)
					log.Printf("Link #%d Text: '%s'\n", index, bookTitle)
					log.Printf("Link #%d Img: '%s'\n", index, bookImg)
					log.Printf("Link #%d Info: '%s'\n", index, bookInfo)
					if index > 5 {
						return false
					}
					return true
				})

				//Preprocessd Flex message Json data

				container := &linebot.CarouselContainer{Type: "carousel", Contents: bubbleContainers}

				flexMessage := linebot.NewFlexMessage("alt text", container)

				// Reply Message

				if _, err = bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
					log.Print(err)
				}

			}
		}
	}
}

func newBubbleContainer(bookTitle, bookLink, bookImg, bookInfo string) *linebot.BubbleContainer {
	return &linebot.BubbleContainer{Type: "bubble",
		Hero:   &linebot.ImageComponent{Type: "image", Size: "full", AspectRatio: "20:13", AspectMode: "fit", URL: bookImg},
		Body:   &linebot.BoxComponent{Type: "box", Layout: "vertical", Spacing: "sm", Contents: []linebot.FlexComponent{&linebot.TextComponent{Type: "text", Text: bookTitle, Wrap: true, Weight: "bold", Size: "xl"}, &linebot.BoxComponent{Type: "box", Layout: "baseline", Contents: []linebot.FlexComponent{&linebot.TextComponent{Type: "text", Text: bookInfo, Wrap: true, Weight: "bold", Size: "xl"}}}}},
		Footer: &linebot.BoxComponent{Type: "box", Layout: "vertical", Spacing: "sm", Contents: []linebot.FlexComponent{&linebot.ButtonComponent{Type: "button", Style: "primary", Action: &linebot.URIAction{Label: "前往", URI: baseComicURL + bookLink}}, &linebot.ButtonComponent{Type: "button", Style: "primary", Action: &linebot.URIAction{Label: "收藏", URI: baseComicURL + bookLink}}}}}
}
