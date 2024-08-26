package time

import "time"

type Duration = time.Duration
type Time = time.Time
type Location = time.Location
type TimerMode = uint8

const InfiniteTimes = -1

const (
	Layout    = "2006-01-02 15:04:05.000"
	LayoutSec = "2006-01-02 15:04:05"
)

const (
	TimerModeMinHeap = 0
	TimerModeWheel   = 1
)

const (
	Nanosecond  = time.Nanosecond
	Microsecond = time.Microsecond
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	Hour        = time.Hour
	Day         = 24 * Hour
)

const (
	stopped int32 = iota
	starting
	started
	stopping
)

var (
	fastForwardTime Duration
)
