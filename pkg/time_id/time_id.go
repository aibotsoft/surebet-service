package time_id

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type TimeId struct {
	seed uint8
	lock sync.Mutex
	last int64
}

func NewTimeId(seed uint8) *TimeId {
	return &TimeId{seed: seed}
}
func (p *TimeId) GetId() int64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	now := time.Now().UnixNano() / 100
	if now <= p.last {
		now = p.last + 1
	}
	p.last = now
	strId := fmt.Sprintf("%d%d", p.seed, p.last)
	value, _ := strconv.ParseInt(strId, 10, 64)
	return value
}
