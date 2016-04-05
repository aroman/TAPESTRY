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
	DB *mgo.Database
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

type ClusterJSON struct {
	models.Cluster
	RootVideo models.Video   `json:"root_video"`
	Videos    []models.Video `json:"videos"`
}

func getClusters(w http.ResponseWriter, r *http.Request) {
	var clusters []models.Cluster

	err := DB.C("clusters").Find(nil).All(&clusters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var clustersJSON []ClusterJSON

	for _, cluster := range clusters {
		cj := ClusterJSON{Cluster: cluster}

		err := DB.C("videos").FindId(cluster.RootVideoID).One(&cj.RootVideo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = DB.C("videos").Find(bson.M{"cluster_id": cluster.ID}).All(&cj.Videos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		clustersJSON = append(clustersJSON, cj)
	}

	json, err := json.Marshal(clustersJSON)
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

	err := DB.C("clusters").UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"label": label}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("ta-da\n"))
	log.Infof("Sent label for cluster %v to %v at %v", id, label, time.Now())

}

func main() {
	log.SetLevel(log.DebugLevel)

	DB = models.GetDB()

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
