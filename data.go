package main

import (
	"container/list"

	. "bitbucket.org/kvu787/boost/lib/angle"
	"bitbucket.org/kvu787/boost/lib/palette"
	. "bitbucket.org/kvu787/boost/lib/vector"

	sf "bitbucket.org/kvu787/gosfml2"
)

var (
	WINDOW         *sf.RenderWindow = nil
	INPUT          *input_s         = &input_s{false, nil}
	GAME_OBJECTS   *list.List       = list.New()
	ASTEROID_COUNT uint             = 3
)

func init() {
	pushFrontAll(GAME_OBJECTS,
		&camera_s{NewZeroVector()},
		&player_s{
			transform_s{NewZeroVector(), NewZeroVector(), NewZeroVector()},
			circle_s{5, 0, 0, palette.BLUE, palette.WHITE},
		},
		&asteroid_s{
			transform_s{NewCartesian(50, -50), NewPolar(50, NewDegrees(-20)), NewZeroVector()},
			circle_s{20, 0, 0, palette.LIGHT_GRAY, palette.WHITE},
			false,
		},
		&asteroid_s{
			transform_s{NewCartesian(0, -100), NewPolar(70, NewDegrees(160)), NewZeroVector()},
			circle_s{20, 0, 0, palette.GRAY, palette.WHITE},
			false,
		},
	)
}

func pushFrontAll(l *list.List, objects ...interface{}) {
	for _, object := range objects {
		l.PushFront(object)
	}
}
