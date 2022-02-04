package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var (
	Rd       *redis.Client
	upgrader = websocket.Upgrader{
		CheckOrigin:      func(r *http.Request) bool { return true },
		HandshakeTimeout: 3 * time.Second,
	}
	ctx = context.Background() // move this into func ws?
)

type ChInit struct {
	InitChannel string `json:"initChannel"`
}

type Msg struct {
	Channel string `json:"channel"`
	Message M      `json:"message"`
}

type M struct {
	Type         string                 `json:"type"`
	FromUserUuid string                 `json:"fromUserUuid"`
	Content      map[string]interface{} `json:"content"`
}

// should handle more errors
// deadlock condition?
func ws(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	// userUuid := c.GetString("user")
	userUuid := "me"
	log.Println(userUuid)

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
					b := []byte(msg.Payload)
					if m, err := decodeToMsg(&b); err == nil {
						if m.Message.FromUserUuid != userUuid {
							err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
							if err != nil {
								log.Println("websocket write err:", err)
								break loop
							}
						}
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

		// chPrefix := strings.Split(string(msg), ":")[0]
		// ch := chPrefix + "-channel"
		// if string(msg) == "test" || string(msg) == "yeah" {
		// 	room <- ch
		// 	log.Println(ch)
		// 	log.Println(<-start)
		// }

		if ch, err := decodeToChInit(&msg); err == nil && ch.InitChannel != "" {
			room <- ch.InitChannel
			<-start
		} else if m, err := decodeToMsg(&msg); err == nil {
			if err := Rd.Publish(ctx, m.Channel, msg).Err(); err != nil {
				log.Println("redis publish err:", err)
				break
			}
		} else {
			log.Println("unknown message type")
		}

	}

}

func GetWs(c *gin.Context) {
	ws(c.Writer, c.Request, c)
}
