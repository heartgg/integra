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
	var err error
	client, err = db.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// read the yaml config of all exams and modalities
	readModalityExams()

	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the IntegriScan API")
	})

	// /ws?roomID=<RoomID>&modality=<Modality>
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})

	// /scan-exams?patientID=<patientID>&modality=<modality>
	http.HandleFunc("/scan-exams", func(w http.ResponseWriter, r *http.Request) {
		scanExamsHandler(client, pool, w, r)
	})
}
