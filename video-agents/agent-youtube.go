package vidagent

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CMU-Perceptual-Computing-Lab/Wisper/database"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type Agent struct {
	service *youtube.Service
}

type SearchParameters struct {
	Terms     string
	Latitude  float64
	Longitude float64
	Radius    string
	TsBefore  time.Time
	TsAfter   time.Time
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
func (agent Agent) Search(params SearchParameters) ([]string, error) {

	var ids []string

	nextPageToken := ""
	for i := 0; i < 5; i++ {
		call := agent.service.Search.List("id").Type("video")

		if params.Longitude != 0 && params.Latitude != 0 {
			call.Location(fmt.Sprintf("%v,%v", params.Longitude, params.Latitude))
		}

		if params.Radius != "" {
			call.LocationRadius(params.Radius)
		}

		if !params.TsBefore.IsZero() {
			call.PublishedBefore(params.TsBefore.Format(time.RFC3339))
		}

		if !params.TsAfter.IsZero() {
			call.PublishedAfter(params.TsAfter.Format(time.RFC3339))
		}

		if params.Terms != "" {
			call.Q(params.Terms)
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
func (agent Agent) GetVideosFromIds(ids []string) ([]*youtube.Video, error) {
	var videos []*youtube.Video
	nextPageToken := ""
	for {
		call := agent.service.Videos.List("id,recordingDetails,snippet,contentDetails")
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

func (agent Agent) GenParams(video *youtube.Video) SearchParameters {

	vm := Serialize(video, "")

	params := SearchParameters{
		Terms: vm.Title,
	}

	params.TsAfter = vm.PublishedAt.AddDate(0, 0, -1)
	params.TsBefore = vm.PublishedAt.AddDate(0, 0, +1)

	if vm.Longitude != 0 {
		params.Longitude = vm.Longitude
		params.Latitude = vm.Latitude
		params.Radius = "100km"
	}

	return params
}

func Serialize(video *youtube.Video, tag string) database.VideoMetadata {

	ts, _ := time.Parse(time.RFC3339, video.Snippet.PublishedAt)

	vm := database.VideoMetadata{
		Title:        video.Snippet.Title,
		PublishedAt:  ts,
		YoutubeID:    video.Id,
		Duration:     video.ContentDetails.Duration,
		Description:  video.Snippet.Description,
		ThumbnailURL: video.Snippet.Thumbnails.Default.Url,
	}

	if tag != "" {
		vm.Tag = tag
	}

	if video.RecordingDetails != nil {
		if video.RecordingDetails.Location != nil {
			vm.Longitude = video.RecordingDetails.Location.Latitude
			vm.Latitude = video.RecordingDetails.Location.Longitude
		}
	}

	return vm
}
