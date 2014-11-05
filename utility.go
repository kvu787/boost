package main

import (
	"container/list"
)

func listWhere(lst *list.List, tag int) interface{} {
	for e := lst.Front(); e != nil; e = e.Next() {
		if e.Value.(Tagged).Tag() == tag {
			return e.Value
		}
	}
	return nil
}

func listSelect(lst *list.List, tag int) *list.List {
	result := list.New()
	for e := lst.Front(); e != nil; e = e.Next() {
		if e.Value.(Tagged).Tag() == tag {
			result.PushBack(e.Value)
		}
	}
	return result
}
