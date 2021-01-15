package loop_count

import (
	"sync"
)

type LoopCount struct {
	count int
	lock  sync.Mutex
}

func NewLoopCount() *LoopCount {
	return &LoopCount{}
}
func (l *LoopCount) Add() {
	l.lock.Lock()
	l.count += 1
	l.lock.Unlock()
}
func (l *LoopCount) Remove() {
	l.lock.Lock()
	l.count -= 1
	l.lock.Unlock()
}
func (l *LoopCount) Get() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.count
}
