//go:build !debug

package time

import "time"

func Now() Time {
	return time.Now()
}
