package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/heartgg/integri-scan/server/pkg/websocket"
	"google.golang.org/api/iterator"
	"gopkg.in/yaml.v3"
)

var modalityExams map[string][]string

func readModalityExams() {
	yfile, err := ioutil.ReadFile("data/modality.yaml")
	if err != nil {
		log.Fatal(err)
	}
	modalityExams = make(map[string][]string)
	err2 := yaml.Unmarshal(yfile, &modalityExams)
	if err2 != nil {
		log.Fatal(err2)
	}
}

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

// controller for getting a scanned barcode id
func scanExamsHandler(client *firestore.Client, pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	patientID := r.URL.Query().Get("patientID")
	modality := r.URL.Query().Get("modality")

	ctx := context.Background()

	if patientID == "" {
		fmt.Fprint(w, "Parameter 'patientID' is required")
		return
	}

	if modality == "" {
		fmt.Fprint(w, "Parameter 'modality' is required")
		return
	}

	// find the first document matching patient id
	pquery := client.Collection("patients").Where("patient_id", "==", patientID).Limit(1).Documents(ctx)
	defer pquery.Stop()
	dsnap, err := pquery.Next()
	if err == iterator.Done {
		fmt.Fprint(w, "No patient found.")
	}

	var patient Patient
	dsnap.DataTo(&patient)

	pjson, _ := json.Marshal(patient)
	pool.Broadcast <- websocket.Message{Body: string(pjson)}
}
