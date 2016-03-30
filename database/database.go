package database

import (
	"time"

	"gopkg.in/mgo.v2"
)

const (
	mongoURI = "mongodb://bambi:bambi@ds019078.mlab.com:19078/tapestry-sandbox"
)

type VideoMetadata struct {
	YoutubeID    string    `bson:"youtube_id" json:"youtube_id"`
	Title        string    `bson:"title" json:"title"`
	Description  string    `bson:"description" json:"description"`
	Duration     string    `bson:"duration" json:"duration"`
	ThumbnailURL string    `bson:"thumbnail_url" json:"thumbnail_url"`
	IsFlagged    bool      `bson:"is_flagged" json:"is_flagged"`
	IsTrashed    bool      `bson:"is_trashed" json:"is_trashed"`
	IsStarred    bool      `bson:"is_starred" json:"is_starred"`
	PublishedAt  time.Time `bson:"published_at" json:"published_at"`
	Latitude     float64   `bson:"latitutde" json:"latitutde"`
	Longitude    float64   `bson:"longitude" json:"longitude"`
	Tag          string    `bson:"tag" json:"tag"`
}

func GetCollection(name string) (*mgo.Collection, error) {
	session, err := mgo.Dial(mongoURI)

	if err != nil {
		return nil, err
	}

	c := session.DB("tapestry-sandbox").C(name)

	index := mgo.Index{
		Key:        []string{"tag", "youtube_id"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}

	err = c.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	return c, nil
}
