package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Video struct {
	ID           bson.ObjectId `bson:"_id,omitempty" json:"id" `
	YoutubeID    string        `bson:"youtube_id" json:"youtube_id"`
	Title        string        `bson:"title" json:"title"`
	Description  string        `bson:"description" json:"description"`
	Duration     string        `bson:"duration" json:"duration"`
	ThumbnailURL string        `bson:"thumbnail_url" json:"thumbnail_url"`
	PublishedAt  time.Time     `bson:"published_at" json:"published_at"`
	Latitude     float64       `bson:"latitutde" json:"latitutde"`
	Longitude    float64       `bson:"longitude" json:"longitude"`
	ClusterID    bson.ObjectId `bson:"cluster_id" json:"cluster_id"`
}
