package websocket

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/LucaTheHacker/go-haversine"
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

		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))

		case message := <-pool.Broadcast:
			fmt.Printf("Handling message: %v\n", message)
			var received models.BroadCastMessage
			if message.Type == 2 {
				json.Unmarshal([]byte(message.Body), &received)
			}
			// list of the workstations with modalities that match the patient modality
			matchedWorkstations := make([]*Client, 0)
			for client := range pool.Clients {
				if message.Type == 2 {
					if received.ExamsResult.Modality == string(client.Modality) {
						// append to list of workstations with matching modalities
						matchedWorkstations = append(matchedWorkstations, client)
					}
				} else if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
			// no workstations with matching modalities
			if len(matchedWorkstations) == 0 {
				return;
			}
			// send exam result to the only matched workstation
			if len(matchedWorkstations) == 1 {
				if err := matchedWorkstations[0].Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
			// loop through list of matched modalities and send exam result to the one nearest to the patient
			var minDistance float64 = 0;
			var nearestWorkstation *Client = matchedWorkstations[0];
			for _, client := range matchedWorkstations {
				clientCoords := haversine.Coordinates{
					Latitude: client.Latitude, 
					Longitude: client.Longitude,
				}
				patientCoords := haversine.Coordinates{
					Latitude: received.Location.Latitude, 
					Longitude: received.Location.Longitude,
				}
				distance := haversine.Distance(clientCoords, patientCoords).Kilometers()
				if minDistance == 0 || distance < minDistance {
					minDistance = distance
					nearestWorkstation = client
				}
			}
			// send message to nearest workstation
			if err := nearestWorkstation.Conn.WriteJSON(message); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
