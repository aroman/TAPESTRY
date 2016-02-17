package main

import (
	"os"
	"fmt"
	"log"
	"net/http"

	"github.com/aroman/tapestry-server/Godeps/_workspace/src/github.com/labstack/echo"
	mw "github.com/aroman/tapestry-server/Godeps/_workspace/src/github.com/labstack/echo/middleware"
	"github.com/aroman/tapestry-server/Godeps/_workspace/src/github.com/rs/cors"
	"github.com/aroman/tapestry-server/Godeps/_workspace/src/google.golang.org/api/googleapi/transport"
	"github.com/aroman/tapestry-server/Godeps/_workspace/src/google.golang.org/api/youtube/v3"
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

	// Make the API call to YouTube.
	call := service.Search.
		List("id,snippet").
		Q(c.Query("q")).
		PublishedAfter(c.Query("after")).
		PublishedBefore(c.Query("before")).
		MaxResults(10)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
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
