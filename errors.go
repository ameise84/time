package time

import (
	"errors"
)

var (
	ErrorTimerNotStop      = errors.New("the timer service has not stopped")
	ErrorFastForwardBefore = errors.New("the time point to fast forward to is before the current time")
)
