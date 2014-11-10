package main

import (
	"container/list"
	"fmt"
	"math"

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

func polynomial(exp, xScale, yScale float64) func(float64) float64 {
	return func(x float64) float64 {
		k := (yScale) / (math.Pow(xScale, exp))
		return k * (math.Pow(x, exp))
	}
}
