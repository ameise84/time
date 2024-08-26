//go:build debug

package time

import "time"

func Now() time.Time {
	return time.Now().Add(fastForwardTime)
}
