package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type VideoMetadata struct {
	Title     string
	Latitude  float64
	Longitude float64
}

type Agent struct {
	service *youtube.Service
}

type SearchParameters struct {
	terms     string
	latitude  float64
	longitude float64
	radius    string
	tsBefore  time.Time
	tsAfter   time.Time
}

func CreateAgent(key string) (*Agent, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: key},
	}

	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	return &Agent{service: service}, nil
}

// searches YouTube based on the provided parameters
// returns a list of youtube video IDs (results)
func (agent Agent) search(params SearchParameters) ([]string, error) {

	// before := timestamp.AddDate(0, 0, -1).Format(time.RFC3339Nano)
	// after := timestamp.AddDate(0, 0, +1).Format(time.RFC3339Nano)
	coordinates := fmt.Sprintf("%v,%v", params.longitude, params.latitude)

	fmt.Printf("time: %v to %v\n", params.tsBefore, params.tsAfter)
	fmt.Printf("place: (%v), radius: %v\n", coordinates, params.radius)
	fmt.Printf("q: %v\n", params.terms)

	var ids []string

	nextPageToken := ""
	for i := 0; i < 5; i++ {
		call := agent.service.Search.List("id").Type("video")

		if params.longitude != 0 && params.latitude != 0 {
			call.Location(fmt.Sprintf("%v,%v", params.longitude, params.latitude))
		}

		if params.radius != "" {
			call.LocationRadius(params.radius)
		}

		if !params.tsBefore.IsZero() {
			call.PublishedBefore(params.tsBefore.Format(time.RFC3339))
		}

		if !params.tsAfter.IsZero() {
			call.PublishedAfter(params.tsAfter.Format(time.RFC3339))
		}

		if params.terms != "" {
			call.Q(params.terms)
		}

		call.PageToken(nextPageToken)
		response, err := call.Do()

		if err != nil {
			return nil, err
		}

		for _, result := range response.Items {
			ids = append(ids, result.Id.VideoId)
		}

		nextPageToken = response.NextPageToken

		if nextPageToken == "" {
			break
		}
	}

	return ids, nil
}

const developerKey = ""

var (
	service *youtube.Service
)

// func printVideo(video *youtube.Video) {
// 	timestamp, _ := time.Parse(time.RFC3339Nano, video.RecordingDetails.RecordingDate)
//
// 	fmt.Printf("Video([%v%v], [%v], [%v])\n", watchBase, video.Id, video.Snippet.Title, timestamp.Format("02/01/2006"))
// }
//
// func printVideos(videos []*youtube.Video) {
// 	fmt.Printf("displaying info about %v videos:\n", len(videos))
// 	for _, video := range videos {
// 		printVideo(video)
// 	}
// }

func (agent Agent) getVideosFromIds(ids []string) ([]*youtube.Video, error) {
	var videos []*youtube.Video
	nextPageToken := ""
	for {
		call := agent.service.Videos.List("id,recordingDetails,snippet")
		call.Id(strings.Join(ids, ","))
		call.PageToken(nextPageToken)

		response, err := call.Do()

		if err != nil {
			return nil, err
		}

		videos = append(videos, response.Items...)

		nextPageToken = response.NextPageToken

		if nextPageToken == "" {
			break
		}
	}

	return videos, nil
}

func depthSearch(location *youtube.GeoPoint, timestamp time.Time, q string) {
	call := service.Search.List("id,snippet").Type("video")
	call.Location(fmt.Sprintf("%v,%v", location.Latitude, location.Longitude))
	call.LocationRadius("100km")
	before := timestamp.AddDate(0, 0, -1).Format(time.RFC3339Nano)
	after := timestamp.AddDate(0, 0, +1).Format(time.RFC3339Nano)
	// fmt.Printf("time: %v to %v\n", before, after)
	// fmt.Printf("place: (%v,%v)\n", location.Latitude, location.Longitude)
	// fmt.Printf("q: %v\n", q)
	call.PublishedBefore(before)
	call.PublishedAfter(after)
	call.Q(q)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making depthSearch search API call: %v", err)
	}
	// j, _ := response.MarshalJSON()
	// fmt.Printf("%v\n", string(j))

	var ids []string
	for _, result := range response.Items {
		ids = append(ids, result.Id.VideoId)
	}

	// videos, err := getVideosFromIds(ids)
	// if err != nil {
	// 	log.Fatalf("Error making videos API call: %v", err)
	// }
	//
	// printVideos(videos)
}
