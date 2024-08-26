package time

import (
	"github.com/ameise84/heap"
)

type timerMinHeap struct {
	h heap.IDHeap[*timer]
}

func (s *timerMinHeap) add(tr *timer) Time {
	_ = s.h.Push(tr.instID, tr)
	return tr.fireAt
}

func (s *timerMinHeap) modify(tr *timer) Time {
	_ = s.h.Update(tr.instID, tr)
	_, top, _ := s.h.Peek()
	return top.fireAt
}

func (s *timerMinHeap) del(tr *timer) {
	s.h.Remove(tr.instID)
}

func (s *timerMinHeap) peek() (*timer, error) {
	_, top, err := s.h.Peek()
	if err != nil {
		return nil, err
	}
	return top, nil
}

func (s *timerMinHeap) pop() (*timer, error) {
	return s.h.Pop()
}

func (s *timerMinHeap) update(tr *timer) {
	s.h.Update(tr.instID, tr)
}
