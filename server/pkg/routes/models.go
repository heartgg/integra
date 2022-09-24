package routes

import (
	"time"
)

type Patient struct {
	ID        string    `firestore:"patient_id" json:"id"`
	Name      string    `firestore:"name" json:"name"`
	Birthdate time.Time `firestore:"birthdate" json:"birthdate"`
	Diagnosis string    `firestore:"diagnosis" json:"diagnosis"`
	Sex       bool      `firestore:"sex" json:"sex"`
}
