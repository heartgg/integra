package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/heartgg/integri-scan/server/models"
	"github.com/heartgg/integri-scan/server/pkg/utils"
	"github.com/heartgg/integri-scan/server/pkg/websocket"
	"google.golang.org/api/iterator"
	"gopkg.in/yaml.v3"
)

var modalityExams map[string][]string
var modalityExamsStr map[string]string

func readModalityExams() {
	yfile, err := ioutil.ReadFile("data/modality.yaml")
	if err != nil {
		log.Fatal(err)
	}
	modalityExams = make(map[string][]string)
	err = yaml.Unmarshal(yfile, &modalityExams)
	if err != nil {
		log.Fatal(err)
	}
	modalityExamsStr = make(map[string]string)
	for key, val := range modalityExams {
		modalityExamsStr[key] = strings.Join(val[:], ", ")
	}
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")

	modality := websocket.ParseModality(r.URL.Query().Get("modality"))
	roomID := r.URL.Query().Get("roomID")
	latitudeStr := r.URL.Query().Get("latitude")
	longitudeStr := r.URL.Query().Get("longitude")

	if modality == "" {
		fmt.Fprint(w, "Parameter 'modality' is required")
		return
	}
	if roomID == "" {
		fmt.Fprint(w, "Parameter 'roomID' is required")
		return
	}
	if latitudeStr == "" {
		fmt.Fprint(w, "Parameter 'latitude' is required")
		return
	}
	if longitudeStr == "" {
		fmt.Fprint(w, "Parameter 'longitude' is required")
		return
	}

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		fmt.Fprint(w, "Parameter 'latitude' is invalid. Must be a float number")
		return
	}
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		fmt.Fprint(w, "Parameter 'latitude' is invalid. Must be a float number")
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
		Latitude: latitude,
		Longitude: longitude,
	}

	pool.Register <- client
	client.Read()
}

// controller for getting a scanned barcode id
func scanExamsHandler(client *firestore.Client, pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	modality := r.URL.Query().Get("modality")
	patientID := r.URL.Query().Get("patientID")
	latitudeStr := r.URL.Query().Get("latitude")
	longitudeStr := r.URL.Query().Get("longitude")

	ctx := context.Background()

	if modality == "" {
		fmt.Fprint(w, "Parameter 'modality' is required")
		return
	}
	if patientID == "" {
		fmt.Fprint(w, "Parameter 'patientID' is required")
		return
	}
	if latitudeStr == "" {
		fmt.Fprint(w, "Parameter 'latitude' is required")
		return
	}
	if longitudeStr == "" {
		fmt.Fprint(w, "Parameter 'longitude' is required")
		return
	}

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		fmt.Fprint(w, "Parameter 'latitude' is invalid. Must be a float number")
		return
	}
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		fmt.Fprint(w, "Parameter 'latitude' is invalid. Must be a float number")
		return
	}

	fmt.Println(patientID, modality)

	// find the first document matching patient id
	pquery := client.Collection("patients").Where("patient_id", "==", patientID).Limit(1).Documents(ctx)
	defer pquery.Stop()
	dsnap, err := pquery.Next()
	if err == iterator.Done {
		fmt.Fprint(w, "No patient found.")
		return
	}

	var patient models.Patient
	err = dsnap.DataTo(&patient)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	exams, err := utils.AskAI(patient.Diagnosis, modalityExams[modality], modalityExamsStr[modality])
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	combined := models.BroadCastMessage{
		ExamsResult: models.ExamsResult{
			Patient:  patient,
			Exams:    exams,
			Modality: modality,
		},
		Location: models.Location{
			Latitude:  latitude,
			Longitude: longitude,
		},
	}

	combinedJson, err := json.Marshal(combined)
	if err != nil {
		fmt.Fprint(w, "Error marshalling combined json")
		return
	}

	pool.Broadcast <- websocket.Message{Type: 2, Body: string(combinedJson)}
}
