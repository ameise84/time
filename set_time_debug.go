//go:build debug

package time

import "time"

func FastForward(dur Duration) error {
	return setFastForwardTime(fastForwardTime + dur)
}

func FastForwardTo(str string, local *Location) error {
	if t, err := time.ParseInLocation(LayoutSec, str, local); err != nil {
		return err
	} else {
		return setFastForwardTime(t.Sub(Now()))
	}
}

func FastForwardToLocal(str string) error {
	return FastForwardTo(str, time.Local)
}

func FastForwardToUTC(str string) error {
	return FastForwardTo(str, time.UTC)
}

func FastForwardToTimeStamp(timestamp int64) error {
	t := time.Unix(timestamp, 0)
	return setFastForwardTime(t.Sub(Now()))
}
