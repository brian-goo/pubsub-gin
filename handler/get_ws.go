package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var (
	Rd       *redis.Client
	upgrader = websocket.Upgrader{}
	ctx      = context.Background()
)

// should handle more errors
// deadlock condition?
func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket connection err:", err)
		return
	}
	defer conn.Close()

	room := make(chan string)
	start := make(chan string)

	go func() {
	loop:
		for {
			sub := Rd.Subscribe(ctx)
			// subCh := sub.Channel()
			defer sub.Close()

			channels := []string{}

			for {
				select {
				case channel := <-room:
					log.Println("channel", channel)
					channels = append(channels, channel)
					sub = Rd.Subscribe(ctx, channels...)
					log.Println("channels", channels)
					// _, err := sub.Receive(ctx)
					// if err != nil {
					// 	log.Println("redis sub connection err:", err)
					// 	break loop
					// }
					start <- "ok"
				case msg := <-sub.Channel():
					log.Println("msg", msg)
					err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
					if err != nil {
						log.Println("websocket write err:", err)
						break loop
					}
				}
			}

		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("websocket read err:", err)
			break
		}
		log.Println(string(msg))

		chPrefix := strings.Split(string(msg), ":")[0]
		ch := chPrefix + "-channel"
		if string(msg) == "test" || string(msg) == "yeah" {
			room <- ch
			log.Println(ch)
			log.Println(<-start)
		}

		if err := Rd.Publish(ctx, ch, msg).Err(); err != nil {
			log.Println("redis publish err:", err)
			break
		}
	}

}

func GetWs(c *gin.Context) {
	ws(c.Writer, c.Request)
}
