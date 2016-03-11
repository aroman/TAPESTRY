package main

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var (
	c *mgo.Collection
)

func getAllVideos(w http.ResponseWriter, r *http.Request) {
	var err error

	var videos []VideoMetadata

	err = c.Find(nil).All(&videos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(videos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(json)
	log.Infof("Sent videos to client at %v", time.Now())
}

func main() {
	log.SetLevel(log.DebugLevel)

	log.Debug("Connecting to database")

	var err error
	c, err = GetCollection("videos")

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", getAllVideos)

	log.Debug("Serving at http://localhost:8000/")
	http.ListenAndServe(":8000", nil)
}
