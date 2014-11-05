package main

import (
	"container/list"
)

// player := listWhere(GAME_OBJECTS, func(i interface{}) bool {
// 	_, ok := i.(player_s)
// 	return ok
// }).(player_s)

func listWhere(lst *list.List, f func(interface{}) bool) interface{} {
	for e := lst.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			return e.Value
		}
	}
	return nil
}

func listSelect(lst *list.List, f func(interface{}) bool) *list.List {
	result := list.New()
	for e := lst.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			result.PushBack(e.Value)
		}
	}
	return result
}
