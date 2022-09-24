package main

import (
	"fmt"
	"net/http"
	
	"github.com/heartgg/integri-scan/server/pkg/websocket"
)


func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")	

	roomID := r.URL.Query().Get("roomID")
	modality := websocket.ParseModality(r.URL.Query().Get("modality"))

	if roomID == "" {
		fmt.Fprint(w, "Parameter 'roomID' is required")
		return
	}
	if modality == "" {
		fmt.Fprint(w, "Parameter 'modality' is required")
		return
	}

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		RoomID: roomID,
		Modality: modality,
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}


func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the IntegriScan API")
	})

	// /ws?roomID=<RoomID>&modality=<Modality>
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	setupRoutes()
	fmt.Println("IntegriScan websocket server is listening on port 8080")
	http.ListenAndServe(":8080", nil)
}