package domain

import "time"

type Patient struct {
	ID               int64     `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Species          string    `json:"species" db:"species"`
	Breed            string    `json:"breed" db:"breed"`
	AproxDateOfBirth time.Time `json:"aproxDateOfBirth" db:"aprox_date_of_birth"`
	OwnerID          int64     `json:"ownerId" db:"owner_id"`
}

func NewPatient(name string, species string, breed string, aproxDateOfBirth time.Time, ownerID int64) *Patient {
	return &Patient{
		Name:             name,
		Species:          species,
		Breed:            breed,
		AproxDateOfBirth: aproxDateOfBirth,
		OwnerID:          ownerID,
	}
}
