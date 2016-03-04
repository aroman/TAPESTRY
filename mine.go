package main

import (
	"bytes"
	"fmt"
	"time"

	"google.golang.org/api/youtube/v3"
	"gopkg.in/alecthomas/kingpin.v2"
)

// type VideoMetadata struct {
// 	Title     string
// 	Latitude  float64
// 	Longitude float64
// }

var (
	agent     *Agent
	latitude  = kingpin.Flag("lat", "latitude of recording").Float64()
	longitude = kingpin.Flag("long", "longitude of recording").Float64()
	radius    = kingpin.Flag("radius", "radius of recording").String()
	before    = kingpin.Flag("before", "uploaded before").String()
	after     = kingpin.Flag("after", "uploaded after").String()
	q         = kingpin.Arg("q", "search query").String()
)

const watchBase = "https://www.youtube.com/watch?v="

func printVideo(video *youtube.Video) {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Video(id=%v", video.Id))

	if video.RecordingDetails != nil {

		if video.RecordingDetails.RecordingDate != "" {
			ts, err := time.Parse(time.RFC3339, video.RecordingDetails.RecordingDate)

			if err != nil {
				panic(err)
			}

			buffer.WriteString(fmt.Sprintf(" date=%v", ts.Format("02/01/2016")))
		}

		if video.RecordingDetails.Location != nil {
			buffer.WriteString(fmt.Sprintf(" lat=%v", video.RecordingDetails.Location.Latitude))
			buffer.WriteString(fmt.Sprintf(" long=%v", video.RecordingDetails.Location.Longitude))
		}

	}

	buffer.WriteString(fmt.Sprintf(" title='%v'", video.Snippet.Title))
	fmt.Printf("%v)\n", buffer.String())
}

func printVideos(videos []*youtube.Video) {
	fmt.Printf("Displaying %v videos:\n", len(videos))
	for _, video := range videos {
		printVideo(video)
	}
}

func main() {
	// session, err := mgo.Dial("mongodb://bambi:bambi@ds019078.mlab.com:19078/tapestry-sandbox")
	//
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()

	// c := session.DB("tapestry-sandbox").C("videos")
	// err = c.Insert(&VideoMetadata{"some title", 0, 0})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	kingpin.Parse()

	var err error
	agent, err := CreateAgent("AIzaSyB-BZx063pUet0zDunRitL_kjwma68tU1c")

	if err != nil {
		panic(err)
	}

	params := SearchParameters{
		terms:     *q,
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

	ids, err := agent.search(params)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", ids)

	videos, err := agent.getVideosFromIds(ids)
	if err != nil {
		panic(err)
	}

	printVideos(videos)
}
