package main

import (
	"encoding/json"
	"net/http"

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
}

func main() {
	session, err := mgo.Dial("mongodb://bambi:bambi@ds019078.mlab.com:19078/tapestry-sandbox")

	if err != nil {
		panic(err)
	}

	defer session.Close()

	c = session.DB("tapestry-sandbox").C("videos")

	http.HandleFunc("/", getAllVideos)
	http.ListenAndServe(":8000", nil)
}
