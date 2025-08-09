package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

var broadcast = make(chan string)

type Broadcaster struct {
	conn *Connector
}

func NewBroadcaster(c *Connector) Broadcaster {
	return Broadcaster{
		conn: c,
	}
}

func (b Broadcaster) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				err := b.conn.closeAll(ctx)
				if err != nil {
					log.Printf("Error closing all clients: %v", err)
					return
				}
				log.Println("Broadcaster stopped by context cancellation")
				return
			case msg := <-broadcast:
				for idx := range b.conn.Clients {
					if b.conn.Clients[idx] == nil {
						log.Printf("Client %d is nil, skipping", idx)
						continue
					}
					err := b.conn.Clients[idx].WriteMessage(websocket.TextMessage, []byte(msg))
					if err != nil {
						log.Printf("Error sending message to client %d: %v", idx, err)
						if err := b.conn.CloseClient(ctx, idx); err != nil {
							log.Printf("Error closing client %d: %v", idx, err)
							return
						}
						log.Printf("Client %d removed due to error", idx)
					}
				}
			}
		}
	}()
}

func (b Broadcaster) NotifyAll(message string) {
	msgJon, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}
	broadcast <- string(msgJon)
}
