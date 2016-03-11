package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"google.golang.org/api/youtube/v3"
	"gopkg.in/alecthomas/kingpin.v2"

	"gopkg.in/mgo.v2"
)

var (
	agent     *Agent
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
	vm := serialize(video, "")

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
	c, err := GetCollection("videos")

	if err != nil {
		panic(err)
	}

	log.Debug("Creating YouTube Agent")
	agent, err := CreateAgent("AIzaSyB-BZx063pUet0zDunRitL_kjwma68tU1c")

	if err != nil {
		panic(err)
	}

	params := SearchParameters{
		terms:     *terms,
		latitude:  *latitude,
		longitude: *longitude,
		radius:    *radius,
	}

	if *before != "" {
		params.tsBefore, err = time.Parse("01-02-2006", *before)
		if err != nil {
			panic(err)
		}
	}

	if *after != "" {
		params.tsAfter, err = time.Parse("01-02-2006", *after)
		if err != nil {
			panic(err)
		}
	}

	log.Debug("Performing root search")
	ids, err := agent.search(params)

	if err != nil {
		panic(err)
	}

	roots, err := agent.getVideosFromIds(ids)
	if err != nil {
		panic(err)
	}

	if len(roots) == 0 {
		log.Fatal("No results found for root search.")
	}
	root := roots[0]

	ids, err = agent.search(agent.genParams(root))
	if err != nil {
		panic(err)
	}

	log.Info("Root video identified")
	printVideo(root)

	videos, err := agent.getVideosFromIds(ids)
	// first result is the root video. skip it.
	videos = videos[1:]

	log.WithFields(log.Fields{
		"results": len(videos),
	}).Info("Found related videos")

	for _, video := range videos {
		printVideo(video)
		// serialize tag is root video's id
		if *dryRun {
			log.Fatal("Dry-run mode; not writing videos to database")
			continue
		}
		err = c.Insert(serialize(video, root.Id))
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
