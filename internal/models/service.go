package models

import "time"

type Service struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Versions    []string  `bson:"versions" json:"versions"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}
