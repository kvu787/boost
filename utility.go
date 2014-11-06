package main

import (
	"container/list"
	"fmt"

	. "bitbucket.org/kvu787/boost/lib/vector"
)

func listAny(l *list.List, f func(interface{}) bool) bool {
	for e := l.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			return true
		}
	}
	return false
}

func listNew(elements ...interface{}) *list.List {
	result := list.New()
	for _, e := range elements {
		fmt.Println(e)
		result.PushBack(e)
	}
	return result
}

func pushFrontAll(l *list.List, objects ...interface{}) {
	for _, object := range objects {
		l.PushFront(object)
	}
}

func getFramedPosition(camera, x Vector) Vector {
	frameTopLeftCorner := camera.Add(NewCartesian(-0.5*float64(WINDOW_SIZE_X), -0.5*float64(WINDOW_SIZE_Y)))
	return x.Sub(frameTopLeftCorner)
}
