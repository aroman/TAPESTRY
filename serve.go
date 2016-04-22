package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CMU-Perceptual-Computing-Lab/Wisper/models"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	db *mgo.Database
)

type Cluster struct {
	models.Cluster
	Videos []models.Video `json:"videos"`
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func getClusters(w http.ResponseWriter, r *http.Request) {
	// var clusters []Cluster

	result := []bson.M{}

	pipe := db.C("clusters").
		Pipe([]bson.M{{
			"$lookup": bson.M{
				"from":         "videos",
				"localField":   "_id",
				"foreignField": "cluster_id",
				"as":           "videos",
			},
		}})

	err := pipe.All(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
	log.Infof("Sent clusters to client at %v", time.Now())
}

func setCluster(w http.ResponseWriter, r *http.Request) {
	label := r.URL.Query().Get("label")
	id := r.URL.Query().Get("id")
	if !bson.IsObjectIdHex(id) {
		http.Error(w, fmt.Sprintf("not an ObjectId: \"%v\"", id), http.StatusBadRequest)
		return
	}
	if !(label == "" || label == "flag" || label == "star" || label == "trash") {
		http.Error(w, fmt.Sprintf("invalid label: \"%v\"", label), http.StatusBadRequest)
		return
	}

	err := db.C("clusters").UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"label": label}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("ta-da\n"))
	log.Infof("Sent label for cluster %v to %v at %v", id, label, time.Now())
}

func main() {
	log.SetLevel(log.DebugLevel)

	db = models.GetDB()

	http.HandleFunc("/", index)
	http.HandleFunc("/api/clusters", getClusters)
	http.HandleFunc("/api/cluster", setCluster)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Debug("Serving at http://localhost:8000/")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		panic(err)
	}
}
