package main

import (
	"bitbucket.org/kvu787/boost/lib/palette"
	. "bitbucket.org/kvu787/boost/lib/vector"
	"container/list"
)

var (
	INPUT        *input_s   = &input_s{false, nil}
	GAME_OBJECTS *list.List = list.New()
)

func init() {
	pushFrontAll(GAME_OBJECTS,
		&player_s{
			transform_s{NewZeroVector(), NewZeroVector(), NewZeroVector()},
			circle_s{10, 0, 0, palette.BLUE, palette.WHITE},
		},
		&asteroid_s{
			transform_s{NewCartesian(50, -60), NewCartesian(-5, 5), NewZeroVector()},
			circle_s{20, 2, 0, palette.GRAY, palette.WHITE},
		}, &asteroid_s{
			transform_s{NewCartesian(-60, 20), NewZeroVector(), NewZeroVector()},
			circle_s{20, 2, 0, palette.GRAY, palette.WHITE},
		},
	)
}

// func newList(elements ...interface{}) *list.List {
// 	result := list.New()
// 	pushFrontAll(result, elements)
// 	return result
// }

func pushFrontAll(l *list.List, objects ...interface{}) {
	for _, object := range objects {
		l.PushFront(object)
	}
}
