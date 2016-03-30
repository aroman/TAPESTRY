package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CMU-Perceptual-Computing-Lab/Wisper/database"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var (
	c *mgo.Collection
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func getAllVideos(w http.ResponseWriter, r *http.Request) {
	var err error

	var videos []database.VideoMetadata

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
	w.Write(json)
	log.Infof("Sent videos to client at %v", time.Now())
}

func main() {
	log.SetLevel(log.DebugLevel)

	log.Debug("Connecting to database")

	var err error
	c, err = database.GetCollection("videos")

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/api/videos", getAllVideos)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Debug("Serving at http://localhost:8000/")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}
