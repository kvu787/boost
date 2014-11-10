package main

import (
	"container/list"
	"time"

	. "bitbucket.org/kvu787/boost/lib/vector"
	. "bitbucket.org/kvu787/boost/objects"

	sf "bitbucket.org/kvu787/gosfml2"
)

const (
	WINDOW_SIZE_X   uint    = 720
	WINDOW_SIZE_Y   uint    = 720
	FPS             uint    = 65
	PLAYER_BOUNDARY float64 = 1000
	SPAWN_BOUNDARY  float64 = 1400
	CAMERA_OFFSET_X float64 = 0
	CAMERA_OFFSET_Y float64 = 0

	SHOULD_SPAWN_ASTEROIDS        bool    = true
	ASTEROID_SPAWN_FREQUENCY      uint    = 3 // lower means more frequent
	ASTEROID_MIN_RADIUS           float64 = 30
	ASTEROID_MAX_RADIUS           float64 = 100
	ASTEROID_MIN_VELOCITY         float64 = 10
	ASTEROID_MAX_VELOCITY         float64 = 400
	ASTEROID_BOUNCE_BACK_VELOCITY float64 = 50
	ASTEROID_INITIAL_SPAWN_COUNT  uint    = 50

	LIGHT_SPAWN_DURATION uint    = 133
	SLIP_WIDTH_SCALING   float64 = 4500 // higher means wider
	MAX_BOOST            float64 = 100
	BOOST_REGENERATION   float64 = 100
	BOOST_BURN           float64 = 100
	SLIP_DURATION        float64 = 4

	PLAYER_RESET_DISTANCE   float64 = 3
	PLAYER_RESET_VELOCITY   float64 = 200
	PLAYER_CONTROL_RADIUS   float64 = 200
	PLAYER_MAX_ACCELERATION float64 = 400
)

// variables that should be immutable
var (
	DURATION_PER_FRAME        time.Duration
	ASTEROID_COLORS           []sf.Color
	WINDOW_DIAGNOL_LENGTH     float64
	PLAYER_ACCELERATION_CURVE func(float64) float64
)

var (
	WINDOW        *sf.RenderWindow
	CAMERA_OFFSET Vector
	INPUT         *Input_s
	PLAYER        *Player_s
	ASTEROIDS     *list.List
	SLIPS         *list.List
	CURRENT_BOOST float64
)
