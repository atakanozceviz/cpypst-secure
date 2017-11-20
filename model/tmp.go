package model

import "sync"

type Tmp struct {
	sync.RWMutex
	tmp string
}

func (t *Tmp) Write(s string) {
	t.Lock()
	defer t.Unlock()
	t.tmp = s
}

func (t *Tmp) Read() string {
	t.Lock()
	defer t.Unlock()
	s := t.tmp
	return s
}
