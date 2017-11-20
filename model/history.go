package model

import (
	"sync"
)

type HistItem struct {
	Ip      string
	Content string
	Time    string
}

type History struct {
	sync.RWMutex
	History []HistItem
}

func (h *History) Add(item HistItem) {
	h.Lock()
	defer h.Unlock()
	h.History = append(h.History, item)
}

func (h *History) Remove(i int) {
	h.Lock()
	defer h.Unlock()
	h.History = append(h.History[:i], h.History[i+1:]...)
}
