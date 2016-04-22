package vidagents

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CMU-Perceptual-Computing-Lab/Wisper/models"

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
func (agent Agent) Search(params SearchParameters) ([]models.Video, error) {

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

	return agent.GetVideosFromIds(ids)
}

// takes a list of video IDs and returns Video objects containing metadata
func (agent Agent) GetVideosFromIds(ids []string) ([]models.Video, error) {

	var videos []models.Video

	nextPageToken := ""
	for {
		call := agent.service.Videos.List("id,recordingDetails,snippet,contentDetails")
		call.Id(strings.Join(ids, ","))
		call.PageToken(nextPageToken)

		response, err := call.Do()

		if err != nil {
			return nil, err
		}

		for _, item := range response.Items {
			videos = append(videos, Serialize(item))
		}

		nextPageToken = response.NextPageToken

		if nextPageToken == "" {
			break
		}
	}

	return videos, nil
}

func (agent Agent) SearchRelatedVideos(id string) ([]string, error) {

	var ids []string

	nextPageToken := ""
	for i := 0; i < 5; i++ {
		call := agent.service.Search.List("id").Type("video")

		call.RelatedToVideoId(id)

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

func (agent Agent) GenParams(video models.Video) SearchParameters {

	params := SearchParameters{
		Terms: video.Title,
	}

	params.TsAfter = video.PublishedAt.AddDate(0, 0, -1)
	params.TsBefore = video.PublishedAt.AddDate(0, 0, +1)

	if video.Longitude != 0 {
		params.Longitude = video.Longitude
		params.Latitude = video.Latitude
		params.Radius = "100km"
	}

	return params
}

func Serialize(video *youtube.Video) models.Video {

	ts, _ := time.Parse(time.RFC3339, video.Snippet.PublishedAt)

	v := models.Video{
		Title:        video.Snippet.Title,
		PublishedAt:  ts,
		YoutubeID:    video.Id,
		Duration:     video.ContentDetails.Duration,
		Description:  video.Snippet.Description,
		ThumbnailURL: video.Snippet.Thumbnails.Default.Url,
	}

	if video.RecordingDetails != nil {
		if video.RecordingDetails.Location != nil {
			v.Longitude = video.RecordingDetails.Location.Latitude
			v.Latitude = video.RecordingDetails.Location.Longitude
		}
	}

	return v
}
