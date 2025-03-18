package time

import (
	"container/list"
	"github.com/ameise84/go_pool"
	"github.com/ameise84/lock"
	"runtime"
	"sync/atomic"
	stdTime "time"
)

var (
	_gTimerLooper Looper
)

func init() {
	_gTimerLooper = newLooper(TimerModeMinHeap)
	_ = _gTimerLooper.Start()
}

type Looper interface {
	Start() error
	Stop()
	NewTimer(cb TimerHandler, ctx any, dur Duration) Timer
	NewTimerWithFireTime(cb TimerHandler, ctx any, fireAt Time) Timer
	NewTick(cb TimerHandler, ctx any, dur Duration, repeat int) Timer
	NewTickWithFirstDo(cb TimerHandler, ctx any, fireAt Time, dur Duration, repeat int) Timer
}

func NewLooper(mod ...TimerMode) Looper {
	m := TimerMode(TimerModeMinHeap)
	if mod != nil && len(mod) == 1 {
		m = mod[0]
	}
	return newLooper(m)
}

func newLooper(mod TimerMode) Looper {
	l := &looper{
		invoke: newTimerInvoker(mod),
		tr:     stdTime.NewTimer(0),
	}
	l.runner = go_pool.NewGoRunner(l, "time runner", go_pool.DefaultOptions())
	l.exec = go_pool.NewGoRunner(l, "time exec", go_pool.DefaultOptions().SetSimCount(1).SetBlock(true).SetCacheMode(true, 512))
	return l
}

type looper struct {
	instID       uint64
	state        atomic.Int32
	isRunning    atomic.Bool
	invoke       timerInvoker
	mu           lock.SpinLock
	runner       go_pool.GoRunner
	exec         go_pool.GoRunner
	tr           *stdTime.Timer
	accTime      Duration //已加速的时长
	nextFireAt   Time
	waitKillList list.List
}

func (l *looper) LogFmt() string {
	return "timer looper"
}

func (l *looper) OnPanic(err error) {
	_gLogger.ErrorBean(l, err)
}

func (l *looper) Start() error {
	if !l.state.CompareAndSwap(stopped, starting) {
		return ErrorTimerNotStop
	}
	l.isRunning.Store(true)
	_ = l.runner.AsyncRun(l.loop)
	l.state.Store(started)
	return nil
}

func (l *looper) Stop() {
	for {
		runtime.Gosched()
		if l.state.Load() == stopped {
			break
		}
		if !l.state.CompareAndSwap(started, stopping) {
			continue
		}
		l.isRunning.Store(false)
		l.tr.Reset(0)
		l.runner.Wait()
		l.exec.Wait()
		unregister(l)
		l.state.Store(stopped)
		break
	}
}

func (l *looper) NewTimer(cb TimerHandler, ctx any, dur Duration) Timer {
	fireAt := Now().Add(dur)
	tr := newTimer(l, cb, ctx, fireAt, dur, 1)
	l.add(tr)
	return tr
}

func (l *looper) NewTimerWithFireTime(cb TimerHandler, ctx any, fireAt Time) Timer {
	tr := newTimer(l, cb, ctx, fireAt, 0, 1)
	l.add(tr)
	return tr
}

func (l *looper) NewTick(cb TimerHandler, ctx any, dur Duration, repeat int) Timer {
	fireAt := Now().Add(dur)
	tr := newTimer(l, cb, ctx, fireAt, dur, repeat)
	l.add(tr)
	return tr
}

func (l *looper) NewTickWithFirstDo(cb TimerHandler, ctx any, fireAt Time, dur Duration, repeat int) Timer {
	tr := newTimer(l, cb, ctx, fireAt, dur, repeat)
	l.add(tr)
	return tr
}

func (l *looper) setFastForwardTime(to Duration) {
	if l.accTime != to {
		l.accTime = to
		l.tr.Reset(0)
	}
}

func (l *looper) add(tr *timer) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !tr.isActive {
		return false
	}
	tr.isJoined = true
	l.invoke.add(tr)
	top, _ := l.invoke.peek()
	if top.instID == tr.instID && !tr.fireAt.Equal(l.nextFireAt) {
		l.nextFireAt = tr.fireAt
		d := tr.fireAt.Sub(Now())
		if d < 0 {
			d = 0
		}
		l.tr.Reset(d)
	}
	return true
}

func (l *looper) modify(tr *timer, at Time, dur Duration, rep int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !tr.isActive {
		return false
	}
	tr.modify(at, dur, rep)
	if tr.isJoined {
		l.invoke.modify(tr)
	} else {
		tr.isJoined = true
		l.invoke.add(tr)
	}
	top, _ := l.invoke.peek()
	if top.instID == tr.instID && !tr.fireAt.Equal(l.nextFireAt) {
		l.nextFireAt = tr.fireAt
		d := tr.fireAt.Sub(Now())
		if d < 0 {
			d = 0
		}
		l.tr.Reset(d)
	}
	return true
}

func (l *looper) pause(tr *timer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !tr.isActive {
		return
	}
	if tr.isJoined {
		l.invoke.del(tr)
		tr.isJoined = false
	}
}

func (l *looper) kill(tr *timer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !tr.isActive {
		return
	}

	if tr.isJoined {
		l.invoke.del(tr)
		tr.isJoined = false
	}
	if tr.isInTrigger.Load() {
		l.waitKillList.PushBack(tr)
	} else {
		freeTimer(tr)
	}
}

func (l *looper) loop(...any) {
loopFor:
	for {
		select {
		case t := <-l.tr.C: //这个返回的是系统真实时间
			d := Duration(0)
			nowTime := t.Add(l.accTime)
			l.mu.Lock()
			for {
				tr, err := l.invoke.peek()
				if err != nil {
					break
				}
				d = tr.fireAt.Sub(nowTime)
				if d > 0 {
					l.tr.Reset(d)
					break
				}

				tr.isInTrigger.Store(true)
				doTimes := tr.fire(nowTime)
				if tr.repeat != 0 {
					l.invoke.update(tr)
				} else {
					tr.isJoined = false
					_, _ = l.invoke.pop()
				}
				_ = l.exec.AsyncRun(fire, tr, doTimes, nowTime)
			}

			for {
				e := l.waitKillList.Front()
				if e == nil {
					break
				}
				ktr := e.Value.(*timer)
				if ktr.isInTrigger.Load() {
					break
				}
				freeTimer(ktr)
				l.waitKillList.Remove(e)
			}

			l.mu.Unlock()
			if !l.isRunning.Load() {
				break loopFor
			}
		}
	}

	for {
		waitLen := l.waitKillList.Len()
		if waitLen == 0 {
			break
		}
		for {
			e := l.waitKillList.Front()
			if e == nil {
				break
			}
			ktr := e.Value.(*timer)
			if ktr.isInTrigger.Load() {
				break
			}
			freeTimer(ktr)
			l.waitKillList.Remove(e)
		}
		stdTime.Sleep(Second)
	}
}

func fire(args ...any) {
	tr := args[0].(*timer)
	do := args[1].(int)
	nowTime := args[2].(Time)
	defer tr.isInTrigger.Store(false)
	tr.cb.OnTimer(tr, do, nowTime)
}
