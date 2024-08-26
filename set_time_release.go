//go:build !debug

package time

func FastForward(d Duration) error {
	return nil
}

func FastForwardTo(timeStr string, local *Location) error {
	return nil
}

func FastForwardToLocal(timeStr string) error {
	return nil
}

func FastForwardToUTC(timeStr string) error {
	return nil
}
