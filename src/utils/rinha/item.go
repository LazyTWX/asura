package rinha

import (
	"fmt"
	"math/rand"
)

type Item struct {
	Type    int     `json:"type"`
	Name    string  `json:"name"`
	Level   int     `json:"level"`
	Payload float64 `json:"payload"`
	Effect  int     `json:"effect"`
}

func GetItemByLevel(level int) int {
	itemTypearr := []*Item{}
	for _, item := range Items {
		if item.Level == level {
			itemTypearr = append(itemTypearr, item)
		}
	}
	selected := itemTypearr[rand.Intn(len(itemTypearr))]
	for i, item := range Items {
		if item.Name == selected.Name {
			return i
		}
	}
	return -1
}

func LevelToString(level int) string {
	return map[int]string{
		0: "comum",
		1: "raro",
		2: "epico",
	}[level]
}

func AddItem(level int, items []int) ([]int, string) {
	item := GetItemByLevel(level)
	if !IsIntInList(item, items) {
		items = append(items, item)
	}
	msg := fmt.Sprintf("Ganhou um item %s\n**%s** (use **j!items** para equipar seu item)", LevelToString(level), Items[item].Name)
	return items, msg
}
