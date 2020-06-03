package handler

import (
	"sync"
)

type StatefulLock struct {
	mu sync.Mutex
	id int64
}

func (l *StatefulLock) Take(id int64) {
	l.mu.Lock()
	l.id = id
}

func (l *StatefulLock) Release(id int64) {
	if l.id == id {
		l.id = 0
		l.mu.Unlock()
	}
}
