package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/rs/cors"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

const developerKey = "AIzaSyB-BZx063pUet0zDunRitL_kjwma68tU1c"

var (
	service *youtube.Service
)

func getVideos(ids []string) ([]*youtube.Video, error) {
	call := service.Videos.List("id,recordingDetails,snippet")
	call.Id(strings.Join(ids, ","))

	response, err := call.Do()
	if err != nil {
		log.Printf("Error making search API call: %v", err)
		return nil, err
	}

	return response.Items, nil
}

func search(c *echo.Context) error {
	call := service.Search.List("id").Type("video")

	if c.Query("after") == "" {
		return c.String(400, "'after' parameter is required")
	}
	call.PublishedAfter(c.Query("after"))

	if c.Query("before") == "" {
		return c.String(400, "'before' parameter is required")
	}
	call.PublishedBefore(c.Query("before"))

	if c.Query("q") != "" {
		call.Q(c.Query("q"))
	}

	if c.Query("maxResults") != "" {
		i, _ := strconv.ParseInt(c.Query("maxResults"), 10, 64)
		call.MaxResults(i)
	}

	if c.Query("location") != "" {
		if c.Query("radius") == "" {
			return c.String(400, "'radius' parameter must accompany 'location' parameter")
		}
		call.Location(c.Query("location"))
		call.LocationRadius(c.Query("radius"))
	}

	results, err := call.Do()
	if err != nil {
		log.Printf("Error making search API call: %v", err)
		return c.String(502, fmt.Sprintf("%v", err))
	}

	var ids []string
	for _, result := range results.Items {
		ids = append(ids, result.Id.VideoId)
	}

	videos, err := getVideos(ids)

	return c.JSONIndent(http.StatusOK, videos, "", "    ")
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(cors.Default().Handler)

	// Routes
	e.Get("/search", search)

	// Initialize YouTube client
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	var err error
	service, err = youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	var port = "5000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	fmt.Printf("Starting server on port %v", port)
	e.Run(":" + port)
}
