package main

import (
	"time"
)

type VideoMetadata struct {
	Title         string
	Latitude      float64
	Longitude     float64
	Tag           string
	RecordingDate time.Time
}
