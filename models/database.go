package models

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

func GetDB() *mgo.Database {
	mongoURI := os.Getenv("TAPESTRY_MONGO_URI")

	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	log.Debugf("Connecting to database %v", mongoURI)
	session, err := mgo.Dial(mongoURI)

	if err != nil {
		panic(err)
	}

	DB := session.DB("tapestry-sandbox")

	err = DB.C("videos").EnsureIndex(mgo.Index{
		Key:        []string{"cluster_id", "youtube_id"},
		Unique:     true,
		Background: true,
	})

	if err != nil {
		panic(err)
	}

	return DB
}
