package websocket

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/heartgg/integri-scan/server/models"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	// keep our workstations connected to heroku
	go func() {
		for range time.Tick(time.Second * 4) {
			pool.Broadcast <- Message{Type: 1, Body: "Ping!"}
		}
	}()

	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			// Send message to each user
			// for client := range pool.Clients {
			// 	fmt.Println(client)
			// 	client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			// }
			break

		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			// Send message to each user
			// for client := range pool.Clients {
			// 	client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			// }
			break

		case message := <-pool.Broadcast:
			fmt.Printf("Handling message: %v\n", message)
			var received models.ExamsResult
			if message.Type == 2 {
				json.Unmarshal([]byte(message.Body), &received)
			}
			// get the client of modality
			fmt.Println("Pool is ",pool.Clients);
			for client := range pool.Clients {
				fmt.Println("Client modality is ",client.Modality," and received is ",received.Modality)
				if message.Type == 2 {
					if received.Modality == string(client.Modality) {
						if err := client.Conn.WriteJSON(message); err != nil {
							fmt.Println(err)
							return
						}
					}
				} else if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
