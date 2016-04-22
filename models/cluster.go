package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Cluster struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	SearchTerms string        `bson:"search_terms" json:"search_terms"`
	Label       string        `bson:"label" json:"label"`
	Latitude    float64       `bson:"latitutde" json:"latitutde"`
	Longitude   float64       `bson:"longitude" json:"longitude"`
	OccurredAt  time.Time     `bson:"occurred_at" json:"occurred_at"`
	MinedAt     time.Time     `bson:"mined_at" json:"mined_at"`
}
