package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Modality string

const (
	IE     Modality = "IE"
	Fluoro Modality = "Fluoro"
	XRAY   Modality = "XRAY"
	CT     Modality = "CT"
	IR     Modality = "IR"
	MRI    Modality = "MRI"
	US     Modality = "US"
	Dexa   Modality = "Dexa"
	NucMed Modality = "NucMed"
)

var (
	modalityMap = map[string]Modality{
		"IE":     IE,
		"Fluoro": Fluoro,
		"XRAY":   XRAY,
		"CT":     CT,
		"IR":     IR,
		"MRI":    MRI,
		"US":     US,
		"Dexa":   Dexa,
		"NucMed": NucMed,
	}
)

func ParseModality(modality string) Modality {
	return modalityMap[modality]
}

type Client struct {
	ID       string
	RoomID   string
	Modality Modality
	Conn     *websocket.Conn
	Pool     *Pool
}

type Message struct {
	Type int
	Body string
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Type: messageType, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
