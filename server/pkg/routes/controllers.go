package routes

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
		RoomID:   roomID,
		Modality: modality,
		Conn:     conn,
		Pool:     pool,
	}

	pool.Register <- client
	client.Read()
}
