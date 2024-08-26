package time

import "time"

func Unix(sec int64, nsec int64) Time {
	return time.Unix(sec, nsec)
}

func Sleep(d Duration) {
	time.Sleep(d)
}

func After(d Duration) <-chan Time {
	return time.After(d)
}

func IsSameDay(t1, t2 Time) bool {
	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day()
}

func IsSameMonth(t1, t2 Time) bool {
	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month()
}

// DaysBetween 从每日零点开始计算相差天数
func DaysBetween(before, after Time) int {
	year1, month1, day1 := before.Date()
	year2, month2, day2 := after.Date()
	if year1 != year2 || month1 != month2 || day1 != day2 {
		d1 := time.Date(year1, month1, day1, 0, 0, 0, 0, before.Location())
		d2 := time.Date(year2, month2, day2, 0, 0, 0, 0, after.Location())
		return int(d2.Sub(d1).Hours() / 24)
	}
	return 0
}

func HoursBetween(before, after Time) int {
	return int(after.Sub(before).Hours())
}
