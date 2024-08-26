package time

import (
	"sync"
	"sync/atomic"
)

var (
	_gTimerPool          sync.Pool
	_gTimerInstIDBuilder atomic.Uint64
)

func init() {
	_gTimerPool = sync.Pool{New: func() any {
		return &timer{instID: _gTimerInstIDBuilder.Add(1)}
	}}
}

func takeTimer(svr *looper, cb TimerHandler, ctx any, fireAt Time, dur Duration, repeat int) *timer {
	tr := _gTimerPool.Get().(*timer)
	tr.isActive = true
	tr.svr = svr
	tr.cb = cb
	tr.ctx = ctx
	tr.fireAt = fireAt
	tr.dur = dur
	tr.repeat = repeat
	return tr
}

func freeTimer(tr *timer) {
	if tr.isActive {
		tr.isActive = false
		tr.svr = nil
		tr.cb = nil
		tr.ctx = nil
		tr.repeat = 0
		_gTimerPool.Put(tr)
	}
}
