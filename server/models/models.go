package models

import "time"

type ExamsResult struct {
	Patient  Patient
	Exams    map[string]int
	Modality string
}

type Location struct {
	Latitude  float64
	Longitude float64
}

type BroadCastMessage struct {
	ExamsResult
	Location
}

type Patient struct {
	ID        string    `firestore:"patient_id" json:"id"`
	Name      string    `firestore:"name" json:"name"`
	Birthdate time.Time `firestore:"birthdate" json:"birthdate"`
	Diagnosis string    `firestore:"diagnosis" json:"diagnosis"`
	Sex       bool      `firestore:"sex" json:"sex"`
}
