// Removes videos and clusters with given label. No-op without --for-sure

package main

import (
	"github.com/CMU-Perceptual-Computing-Lab/Wisper/models"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	label   = kingpin.Flag("label", "remove clusters with this label").Required().String()
	forSure = kingpin.Flag("for-sure", "really delete stuff").Bool()
)

func main() {
	log.SetLevel(log.DebugLevel)

	kingpin.Parse()

	DB := models.GetDB()
	log.
		WithField("label", *label).
		Info("Looking for videos and clusters...")

	var trashedClusters []models.Cluster
	var trashedClusterIDs []bson.ObjectId

	DB.C("clusters").
		Find(bson.M{"label": *label}).
		Select(bson.M{"cluster_id": true}).
		All(&trashedClusters)

	log.WithFields(log.Fields{
		"count": len(trashedClusters),
	}).Info("Query for clusters complete")

	if !*forSure {
		log.Fatalf("you didn't pass --for-sure, exiting safely...")
	}

	trashedClusterIDs = make([]bson.ObjectId, len(trashedClusters))
	for i, cluster := range trashedClusters {
		trashedClusterIDs[i] = cluster.ID
	}

	info, err := DB.C("videos").
		RemoveAll(bson.M{"cluster_id": bson.M{"$in": trashedClusterIDs}})

	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{
		"count": info.Removed,
	}).Info("deleted videos")

	info, err = DB.C("clusters").
		RemoveAll(bson.M{"_id": bson.M{"$in": trashedClusterIDs}})

	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{
		"count": info.Removed,
	}).Info("deleted clusters")

}
