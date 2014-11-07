package main

import (
	"container/list"
	"time"

	. "bitbucket.org/kvu787/boost/lib/vector"
	. "bitbucket.org/kvu787/boost/objects"

	sf "bitbucket.org/kvu787/gosfml2"
)

const (
	WINDOW_SIZE_X                 uint    = 720
	WINDOW_SIZE_Y                 uint    = 720
	FPS                           uint    = 65
	ASTEROID_LIMIT                uint    = 40
	PLAYER_BOUNDARY               float64 = 1000
	SPAWN_BOUNDARY                float64 = 1400
	PLAYER_RESET_DISTANCE         float64 = 3
	PLAYER_RESET_VELOCITY         float64 = 200
	ASTEROID_SPAWN_FREQUENCY      uint    = 3 // lower means more frequent
	CAMERA_OFFSET_X               float64 = 0
	CAMERA_OFFSET_Y               float64 = 0
	ASTEROID_MIN_RADIUS           float64 = 30
	ASTEROID_MAX_RADIUS           float64 = 100
	ASTEROID_MIN_VELOCITY         float64 = 10
	ASTEROID_MAX_VELOCITY         float64 = 400
	SENSITIVITY                   float64 = 2
	ASTEROID_BOUNCE_BACK_VELOCITY float64 = 300
	ASTEROID_INITIAL_SPAWN_COUNT  uint    = 100
	SHOULD_SPAWN_ASTEROIDS        bool    = true
	LIGHT_SPAWN_DURATION          uint    = 1337
	SLIP_WIDTH_SCALING            float64 = 45000 // higher means wider
	MAX_BOOST                     float64 = 100
	BOOST_REGENERATION            float64 = 100
	BOOST_BURN                    float64 = 100
	SLIP_DURATION                 float64 = 4
)

var (
	WINDOW                *sf.RenderWindow
	DURATION_PER_FRAME    time.Duration
	CAMERA_OFFSET         Vector
	INPUT                 *Input_s
	PLAYER                *Player_s
	ASTEROIDS             *list.List
	ASTEROID_COLORS       []sf.Color
	SLIPS                 *list.List
	WINDOW_DIAGNOL_LENGTH float64
	CURRENT_BOOST         float64
)
