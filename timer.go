package time

import (
	"github.com/ameise84/heap/compare"
	"sync/atomic"
	"time"
)

type TriggerMode int

type TimerHandler interface {
	OnTimer(tr Timer, doTimes int, now Time)
}

type Timer interface {
	Context() any
	Reset(dur Duration, repeat ...int) (Time, bool) //重置成功,返回下一次执行的时间,如果原来的定时任务还存在,原任务将会被kill(如果已经被触发,还是会触发),
	ResetFireAt(fire Time) bool
	Kill()  //kill以后,对象不能在使用,应用层不可继续操作该对象
	Pause() //暂停后可继续持有对象Reset
}

const (
	TriggerModeOnce TriggerMode = 0
	TriggerModeMore TriggerMode = 1
)

var (
	zeroTime Time
)

func NewTimer(cb TimerHandler, ctx any, dur Duration) Timer {
	return _gTimerLooper.NewTimer(cb, ctx, dur)
}

func NewTimerWithFireTime(cb TimerHandler, ctx any, fireAt Time) Timer {
	return _gTimerLooper.NewTimerWithFireTime(cb, ctx, fireAt)
}

func NewTick(cb TimerHandler, ctx any, dur Duration, repeat int) Timer {
	return _gTimerLooper.NewTick(cb, ctx, dur, repeat)
}

func NewTickWithFirstDo(cb TimerHandler, ctx any, fireAt Time, dur Duration, repeat int) Timer {
	return _gTimerLooper.NewTickWithFirstDo(cb, ctx, fireAt, dur, repeat)
}

func newTimer(lp *looper, cb TimerHandler, ctx any, fireAt Time, dur Duration, repeat int) *timer {
	return takeTimer(lp, cb, ctx, fireAt, dur, repeat)
}

type timer struct {
	isWaitKill  atomic.Bool
	isInTrigger atomic.Bool
	instID      uint64
	isActive    bool
	isJoined    bool
	svr         *looper
	ctx         any
	cb          TimerHandler
	fireAt      time.Time
	dur         Duration
	repeat      int //重复次数
	mode        TriggerMode
}

func (tr *timer) modify(fireAt Time, dur Duration, repeat int) {
	tr.fireAt = fireAt
	tr.dur = dur
	tr.repeat = repeat
}

func (tr *timer) fire(now time.Time) int {
	doTimes := 1
	if tr.mode == TriggerModeMore && tr.repeat != 1 {
		doTimes += int(now.Sub(tr.fireAt) / tr.dur)
		if tr.repeat != InfiniteTimes && doTimes > tr.repeat {
			doTimes = tr.repeat
		}
	}

	if tr.repeat != InfiniteTimes {
		tr.repeat -= doTimes
	}

	if tr.repeat != 0 {
		tr.fireAt = tr.fireAt.Add(Duration(doTimes) * tr.dur)
	}
	return doTimes
}

func (tr *timer) Compare(c compare.Ordered) compare.Result {
	ctr := c.(*timer)
	d := tr.fireAt.Sub(ctr.fireAt)
	if d < 0 {
		return compare.Smaller
	}
	if d > 0 {
		return compare.Larger
	}
	return compare.Equal
}

func (tr *timer) Context() any {
	return tr.ctx
}

func (tr *timer) Reset(dur Duration, repeat ...int) (Time, bool) {
	rep := 1
	if !tr.isActive {
		return zeroTime, false
	}
	if repeat != nil {
		rep = repeat[0]
		if !(rep == InfiniteTimes || rep > 0) {
			return zeroTime, false
		}
	}
	if dur == 0 && rep != 1 {
		return zeroTime, false
	}
	now := Now()
	fireAt := now.Add(dur)
	tr.svr.modify(tr, fireAt, dur, rep)
	return fireAt, true
}

func (tr *timer) ResetFireAt(fireAt Time) bool {
	if !tr.isActive {
		return false
	}
	tr.svr.modify(tr, fireAt, tr.dur, 1)
	return true
}

func (tr *timer) Kill() {
	if !tr.isActive {
		return
	}
	tr.svr.kill(tr)
}

func (tr *timer) Pause() {
	if !tr.isActive {
		return
	}
	tr.svr.pause(tr)
}
