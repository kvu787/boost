package main

import (
	. "bitbucket.org/kvu787/boost/objects"
	"container/list"
	"time"

	. "bitbucket.org/kvu787/boost/lib/vector"

	sf "bitbucket.org/kvu787/gosfml2"
)

const (
	WINDOW_SIZE_X                 uint    = 1000
	WINDOW_SIZE_Y                 uint    = 720
	FPS                           uint    = 65
	ASTEROID_LIMIT                uint    = 40
	PLAYER_BOUNDARY               float64 = 1000
	SPAWN_BOUNDARY                float64 = 1200
	PLAYER_RESET_DISTANCE         float64 = 3
	PLAYER_RESET_VELOCITY         float64 = 200
	ASTEROID_SPAWN_FREQUENCY      uint    = 20
	CAMERA_OFFSET_X               float64 = 0
	CAMERA_OFFSET_Y               float64 = 0
	ASTEROID_MIN_RADIUS           float64 = 30
	ASTEROID_MAX_RADIUS           float64 = 100
	ASTEROID_MIN_VELOCITY         float64 = 10
	ASTEROID_MAX_VELOCITY         float64 = 300
	SENSITIVITY                   float64 = 2
	ASTEROID_BOUNCE_BACK_VELOCITY float64 = 300
	ASTEROID_INITIAL_SPAWN_COUNT  uint    = 50
)

var (
	WINDOW             *sf.RenderWindow
	DURATION_PER_FRAME time.Duration
	CAMERA_OFFSET      Vector
	INPUT              *Input_s
	PLAYER             *Player_s
	ASTEROIDS          *list.List
	ASTEROID_COLORS    []sf.Color
)
