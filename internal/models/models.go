package models

import (
	"sync"
)

type NumberSystem struct {
	CurrentNumbers []int `json:"currentNumbers"`
	WaitingNumbers []int `json:"waitingNumbers"`
	Mutex          sync.RWMutex `json:"-"`
}

type Category struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	RangeStart int    `json:"rangeStart"`
	RangeEnd   int    `json:"rangeEnd"`
	Color      string `json:"color"`
}

type Config struct {
	Categories []Category `json:"categories"`
}
