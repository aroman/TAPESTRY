// Removes trashed videos and clusters. Supports the --dry flag.

package main

import (
	"github.com/CMU-Perceptual-Computing-Lab/Wisper/models"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	dryRun = kingpin.Flag("dry", "dry run").Bool()
)

func main() {
	log.SetLevel(log.DebugLevel)

	kingpin.Parse()

	DB := models.GetDB()
	log.Info("Looking for trashed videos and clusters...")

	var trashedClusters []models.Cluster
	var trashedClusterIDs []bson.ObjectId

	DB.C("clusters").
		Find(bson.M{"label": "trash"}).
		Select(bson.M{"cluster_id": true}).
		All(&trashedClusters)

	log.WithFields(log.Fields{
		"count": len(trashedClusters),
	}).Info("found trashed clusters")

	if *dryRun {
		log.Fatalf("dry-run; exiting")
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
	}).Info("deleted trashed videos")

	info, err = DB.C("clusters").
		RemoveAll(bson.M{"_id": bson.M{"$in": trashedClusterIDs}})

	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{
		"count": info.Removed,
	}).Info("deleted trashed clusters")

}
