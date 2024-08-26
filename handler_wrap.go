package time

type HandlerWrap interface {
	GetTimerHandler() TimerHandler
	C() <-chan Trigger
}

func NewTimerHandlerWrap(size int) HandlerWrap {
	return &handlerWrap{
		c: make(chan Trigger, size),
	}
}

type handlerWrap struct {
	c chan Trigger
}

func (hw *handlerWrap) GetTimerHandler() TimerHandler {
	return hw
}

func (hw *handlerWrap) OnTimer(tr Timer, doTimes int, now Time) {
	hw.c <- Trigger{
		tr:      tr,
		doTimes: doTimes,
		now:     now,
	}
}

func (hw *handlerWrap) C() <-chan Trigger {
	return hw.c
}

type Trigger struct {
	tr      Timer
	doTimes int
	now     Time
}

func (tr Trigger) Timer() Timer {
	return tr.tr
}

func (tr Trigger) DoTimes() int {
	return tr.doTimes
}

func (tr Trigger) Now() Time {
	return tr.now
}
