package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

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

// takes a list of video IDs and returns Video objects containing metadata
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

func (agent Agent) genParams(video *youtube.Video) SearchParameters {

	vm := serialize(video, "")

	params := SearchParameters{
		terms: vm.Title,
	}

	if !vm.RecordingDate.IsZero() {
		params.tsAfter = vm.RecordingDate.AddDate(0, 0, -1)
		params.tsBefore = vm.RecordingDate.AddDate(0, 0, +1)
	}

	if vm.Longitude != 0 {
		params.longitude = vm.Longitude
		params.latitude = vm.Latitude
		params.radius = "100km"
	}

	return params
}

func serialize(video *youtube.Video, tag string) VideoMetadata {

	vm := VideoMetadata{
		Title: video.Snippet.Title,
	}

	if tag != "" {
		vm.Tag = tag
	}

	if video.RecordingDetails != nil {

		if video.RecordingDetails.RecordingDate != "" {
			ts, _ := time.Parse(time.RFC3339, video.RecordingDetails.RecordingDate)
			vm.RecordingDate = ts
		}

		if video.RecordingDetails.Location != nil {
			if video.RecordingDetails.Location.Latitude != 0 {
				vm.Longitude = video.RecordingDetails.Location.Latitude
			}
			if video.RecordingDetails.Location.Latitude != 0 {
				vm.Latitude = video.RecordingDetails.Location.Longitude
			}
		}

	}

	return vm
}
