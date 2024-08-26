package time

import (
	"github.com/ameise84/heap"
)

type timerInvoker interface {
	add(tr *timer) Time
	modify(tr *timer) Time
	del(tr *timer)
	peek() (*timer, error)
	pop() (*timer, error)
	update(tr *timer)
}

func newTimerInvoker(mod TimerMode) timerInvoker {
	switch mod {
	case TimerModeMinHeap:
		return &timerMinHeap{heap.NewIDHeapMin[*timer]()}
	case TimerModeWheel:
		panic("TimerModeWheel not supported")
	default:
		panic("TimerMode unknown")
	}
}
