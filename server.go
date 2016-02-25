package main

import (
	"os"
	"fmt"
	"log"
	"strconv"
	"net/http"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/rs/cors"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

const developerKey = "AIzaSyB-BZx063pUet0zDunRitL_kjwma68tU1c"

func hello(c *echo.Context) error {

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	call := service.Search.List("id,snippet")

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

	if c.Query("radius") != "" {
		call.LocationRadius(c.Query("radius"))
	}

	if c.Query("location") != "" {
		call.Location(c.Query("location"))
	}

	// // Make the API call to YouTube.
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making search API call: %v", err)
		return c.String(502, fmt.Sprintf("%v", err))
	}

	return c.JSONIndent(http.StatusOK, response.Items, "", "    ")
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(cors.Default().Handler)

	// Routes
	e.Get("/search", hello)

	var port string = "5000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	fmt.Printf("Starting server on port %v", port)
	// Start server
	e.Run(":" + port)
}
