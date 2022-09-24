package websocket

import (
	"encoding/json"
	"fmt"
	"time"
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
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}
			break

		case message := <-pool.Broadcast:
			fmt.Printf("Handling message: %v\n", message)
			var received map[string]string
			if message.Type == 2 {
				json.Unmarshal([]byte(message.Body), &received)
			}
			// get the client of modality
			for client := range pool.Clients {
				if message.Type == 2 {
					fmt.Println(received["modality"])
					fmt.Println(client.Modality)
					if received["modality"] == string(client.Modality) {
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
