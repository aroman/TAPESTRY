package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/CMU-Perceptual-Computing-Lab/Wisper/models"
	"github.com/CMU-Perceptual-Computing-Lab/Wisper/video-agents"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	agent     *vidagents.Agent
	latitude  = kingpin.Flag("lat", "latitude of recording").Float64()
	longitude = kingpin.Flag("long", "longitude of recording").Float64()
	radius    = kingpin.Flag("radius", "radius of recording").String()
	before    = kingpin.Flag("before", "uploaded before").String()
	after     = kingpin.Flag("after", "uploaded after").String()
	terms     = kingpin.Flag("terms", "search query").Required().String()
	tag       = kingpin.Flag("tag", "tag videos with tag").String()
	dryRun    = kingpin.Flag("dry", "dry run (don't save videos)").Bool()
)

func truncate(str string, max int) string {
	if max < len(str) {
		return str[:max-3] + "..."
	}
	return str
}

func printVideo(video models.Video) {
	log.WithFields(log.Fields{
		"id":        video.YoutubeID,
		"published": video.PublishedAt.Format("01/02/2006"),
		"lat":       video.Latitude,
		"long":      video.Longitude,
	}).Info(truncate(video.Title, 44))
}

func main() {

	log.SetLevel(log.DebugLevel)

	kingpin.Parse()

	// Check if there's already a cluster with our search terms
	var existingCluster models.Cluster

	DB.C("clusters").Find(bson.M{"search_terms": params.Terms}).One(&existingCluster)

	if existingCluster.ID != "" {
		log.Fatalf("Cluster already mined (there's another cluster with the same search terms)")
	}

	DB := models.GetDB()

	log.Debug("Creating YouTube Agent")
	agent, err := vidagents.CreateAgent("AIzaSyB-BZx063pUet0zDunRitL_kjwma68tU1c")

	if err != nil {
		panic(err)
	}

	params := vidagents.SearchParameters{
		Terms:     *terms,
		Latitude:  *latitude,
		Longitude: *longitude,
		Radius:    *radius,
	}

	if *before != "" {
		params.TsBefore, err = time.Parse("01-02-2006", *before)
		if err != nil {
			panic(err)
		}
	}

	if *after != "" {
		params.TsAfter, err = time.Parse("01-02-2006", *after)
		if err != nil {
			panic(err)
		}
	}

	log.Debugf("Performing root search: %v", *terms)
	roots, err := agent.Search(params)

	if err != nil {
		panic(err)
	}

	if len(roots) == 0 {
		log.Fatal("Root search returned 0 results")
	}

	// Map from YouTube IDs to number of occurances in results
	m := make(map[models.Video]int)

	for _, root := range roots {
		videos, err := agent.Search(vidagents.SearchParameters{
			Terms: root.Title,
		})

		if err != nil {
			panic(err)
		}

		fmt.Printf("Found %v videos\n", len(videos))

		for _, video := range videos {
			m[video]++
		}
	}

	// XXX: Refactor
	var goodVideos []models.Video

	n := map[int][]models.Video{}
	var a []int
	for k, v := range m {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	for _, k := range a {
		for _, video := range n[k] {
			if k >= 3 {
				goodVideos = append(goodVideos, video)
				// fmt.Printf("%d: %v\n", k, video)
			}
		}
	}

	if *dryRun {
		log.Warn("Dry-run mode; not writing to database")
		os.Exit(0)
	}

	cluster := models.Cluster{
		SearchTerms: params.Terms,
		MinedAt:     time.Now(),
	}
	cluster.ID = bson.NewObjectId()

	for _, video := range goodVideos {
		video.ID = bson.NewObjectId()
		video.ClusterID = cluster.ID

		err = DB.C("videos").Insert(video)
		if err != nil {
			panic(err)
		}
	}

	err = DB.C("clusters").Insert(cluster)
	if err != nil {
		panic(err)
	}

	log.Debugf("Cluster of %v video(s) written to database", len(goodVideos))
}
