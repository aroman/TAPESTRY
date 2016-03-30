package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aroman/tapestry/database"
	"github.com/aroman/tapestry/video-agents"
	"google.golang.org/api/youtube/v3"
	"gopkg.in/alecthomas/kingpin.v2"

	"gopkg.in/mgo.v2"
)

var (
	agent     *vidagent.Agent
	latitude  = kingpin.Flag("lat", "latitude of recording").Float64()
	longitude = kingpin.Flag("long", "longitude of recording").Float64()
	radius    = kingpin.Flag("radius", "radius of recording").String()
	before    = kingpin.Flag("before", "uploaded before").String()
	after     = kingpin.Flag("after", "uploaded after").String()
	terms     = kingpin.Flag("terms", "search query").String()
	tag       = kingpin.Flag("tag", "tag videos with tag").String()
	dryRun    = kingpin.Flag("dry", "dry run (don't save videos)").Bool()
)

func truncate(str string, max int) string {
	if max < len(str) {
		return str[:max-3] + "..."
	}
	return str
}

func printVideo(video *youtube.Video) {
	vm := vidagent.Serialize(video, "")

	log.WithFields(log.Fields{
		"id":        vm.YoutubeID,
		"published": vm.PublishedAt.Format("01/02/2006"),
		"lat":       vm.Latitude,
		"long":      vm.Longitude,
	}).Info(truncate(vm.Title, 44))
}

func main() {

	log.SetLevel(log.DebugLevel)

	kingpin.Parse()

	var err error

	log.Debug("Connecting to database")
	c, err := database.GetCollection("videos")

	if err != nil {
		panic(err)
	}

	log.Debug("Creating YouTube Agent")
	agent, err := vidagent.CreateAgent("AIzaSyB-BZx063pUet0zDunRitL_kjwma68tU1c")

	if err != nil {
		panic(err)
	}

	params := vidagent.SearchParameters{
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
		log.Fatal("No results found for root search")
	}
	root := roots[0]

	ids, err = agent.Search(agent.GenParams(root))
	if err != nil {
		panic(err)
	}

	log.Info("Root video identified")
	printVideo(root)

	videos, err := agent.GetVideosFromIds(ids)

	log.WithFields(log.Fields{
		"results": len(videos),
	}).Info("Found related videos")

	for _, video := range videos {
		printVideo(video)
		if *dryRun {
			log.Fatal("Dry-run mode; not writing videos to database")
			continue
		}
		// Serialize tag is root video's id
		err = c.Insert(vidagent.Serialize(video, root.Id))
		if err != nil {
			if mgo.IsDup(err) {
				log.Warn("Video already in database; skipping")
				continue
			}
			panic(err)
		}
		log.Debug("Video written to database")
	}

}
