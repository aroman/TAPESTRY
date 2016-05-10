package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/CMU-Perceptual-Computing-Lab/Wisper/models"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	db *mgo.Database
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

// XXX: This is not a great way to mine for clusters, it was added last-minute
// as a minor convenience over running the command-line miner.
func mine(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/mine.html")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}

	terms := r.PostFormValue("terms")
	fmt.Println(terms)

	cmd := exec.Command("./mine-youtube", "--terms", terms)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	go cmd.Run()

	w.Write([]byte("Job running. It may take a few minutes to show up in the interface."))
}

func getClusters(w http.ResponseWriter, r *http.Request) {
	clusters := []bson.M{}

	pipe := db.C("clusters").
		Pipe([]bson.M{{
			"$lookup": bson.M{
				"from":         "videos",
				"localField":   "_id",
				"foreignField": "cluster_id",
				"as":           "videos",
			},
		}})

	err := pipe.All(&clusters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(clusters)
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
	http.HandleFunc("/mine", mine)
	http.HandleFunc("/api/clusters", getClusters)
	http.HandleFunc("/api/cluster", setCluster)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Debug("Serving at http://localhost:8000/")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		panic(err)
	}
}
