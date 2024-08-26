package time

import (
	"github.com/ameise84/lock"
)

var (
	_gAccTimeLocker lock.SpinLock
	looperMap       = make(map[uint64]*looper, 32)
)

func register(l *looper) {
	_gAccTimeLocker.Lock()
	looperMap[l.instID] = l
	l.setFastForwardTime(fastForwardTime)
	_gAccTimeLocker.Unlock()
}

func unregister(l *looper) {
	_gAccTimeLocker.Lock()
	delete(looperMap, l.instID)
	_gAccTimeLocker.Unlock()
}

func setFastForwardTime(d Duration) error {
	_gAccTimeLocker.Lock()
	defer _gAccTimeLocker.Unlock()
	if d < fastForwardTime {
		return ErrorFastForwardBefore
	}
	fastForwardTime = d
	for _, lp := range looperMap {
		lp.setFastForwardTime(d)
	}
	return nil
}
