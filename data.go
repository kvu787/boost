package main

import (
	. "bitbucket.org/kvu787/boost/objects"
	"container/list"
	"time"

	"bitbucket.org/kvu787/boost/lib/palette"
	. "bitbucket.org/kvu787/boost/lib/vector"

	sf "bitbucket.org/kvu787/gosfml2"
)

const (
	WINDOW_SIZE_X         uint    = 1000
	WINDOW_SIZE_Y         uint    = 720
	FPS                   uint    = 65
	ASTEROID_LIMIT        uint    = 40
	PLAYER_BOUNDARY       float64 = 1000
	SPAWN_BOUNDARY        float64 = 1800
	PLAYER_RESET_DISTANCE float64 = 3
)

var (
	WINDOW             *sf.RenderWindow = nil
	DURATION_PER_FRAME time.Duration    = time.Duration(int64(time.Second) / int64(FPS))
	CAMERA_OFFSET      Vector           = NewZeroVector()
	INPUT              *Input_s         = &Input_s{false, nil}

	PLAYER *Player_s = &Player_s{
		Transform_s{
			NewZeroVector(),
			NewZeroVector(),
			NewZeroVector(),
		},
		Circle_s{10},
		RenderProperties_s{0, 0, palette.BLUE, palette.WHITE}}

	ASTEROIDS *list.List = listNew(
		&Asteroid_s{
			Transform_s{
				NewPolar(30, 0),
				NewZeroVector(),
				NewZeroVector(),
			},
			Circle_s{20},
			RenderProperties_s{0, 0, palette.LIGHT_GRAY, palette.WHITE}},
		&Asteroid_s{
			Transform_s{
				NewPolar(150, DegreesToRadians(240)),
				NewPolar(60, DegreesToRadians(60)),
				NewZeroVector(),
			},
			Circle_s{15},
			RenderProperties_s{0, 0, palette.GRAY, palette.WHITE}})
)
