package objects

import (
	sf "bitbucket.org/kvu787/gosfml2"
	"time"
)

type Slip_s struct {
	Segment_s
	Width       float64
	Color       sf.Color
	Duration    time.Duration
	TimeSpawned time.Time
}
