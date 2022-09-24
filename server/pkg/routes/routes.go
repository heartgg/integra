package routes

import (
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/heartgg/integri-scan/server/pkg/db"
	"github.com/heartgg/integri-scan/server/pkg/websocket"
)

var client *firestore.Client

func SetupRoutes() {
	// connect to the DB (client)
	client, err := db.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

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
