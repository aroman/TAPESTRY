package main

import (
	"os"
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
	ids, err := agent.Search(params)

	if err != nil {
		panic(err)
	}

	roots, err := agent.GetVideosFromIds(ids)
	if err != nil {
		panic(err)
	}

	if len(roots) == 0 {
		log.Fatal("Root search returned 0 results")
	}

	root := roots[0]

	ids, err = agent.Search(agent.GenParams(root))
	if err != nil {
		panic(err)
	}

	videos, err := agent.GetVideosFromIds(ids)

	for _, video := range videos {
		printVideo(video)
	}

	if *dryRun {
		log.Warn("Dry-run mode; not writing to database")
		os.Exit(0)
	}

	cluster := models.Cluster{
		SearchTerms: params.Terms,
		MinedAt:     time.Now(),
		Latitude:    root.Latitude,
		Longitude:   root.Longitude,
		// TODO: Don't assume the video's publishing date is the same as the occurance date
		OccurredAt: root.PublishedAt,
	}
	cluster.ID = bson.NewObjectId()

	// Check if there's already a cluster with a root video with our same youtube ID.
	var existingRoot models.Video
	var existingCluster models.Cluster

	DB.C("videos").Find(bson.M{"youtube_id": root.YoutubeID}).One(&existingRoot)
	DB.C("clusters").Find(bson.M{"root_video_id": existingRoot.ID}).One(&existingCluster)

	if existingCluster.ID != "" {
		log.Fatalf("Cluster already mined (there's another cluster with the same root video)")
	}

	for _, video := range videos {
		// Associate the root video with the cluster
		if video.YoutubeID == root.YoutubeID {
			video.ID = bson.NewObjectId()
			cluster.RootVideoID = video.ID
		}

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

	log.Debugf("Cluster of %v video(s) written to database", len(videos))
}
