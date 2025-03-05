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

func ZeroTimeUTC(t Time) Time {
	return t.Truncate(time.Hour * 24).Local()
}

func ZeroTimeLocal(t Time) Time {
	y1, m1, d1 := t.Local().Date()
	return time.Date(y1, m1, d1, 0, 0, 0, 0, time.Local)
}

func IsSameDayUTC(t1, t2 Time) bool {
	y1, m1, d1 := t1.UTC().Date()
	y2, m2, d2 := t2.UTC().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
func IsSameDayLocal(t1, t2 Time) bool {
	y1, m1, d1 := t1.Local().Date()
	y2, m2, d2 := t2.Local().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func IsSameMonthUTC(t1, t2 Time) bool {
	y1, m1, _ := t1.UTC().Date()
	y2, m2, _ := t2.UTC().Date()
	return y1 == y2 && m1 == m2
}

func IsSameMonthLocal(t1, t2 Time) bool {
	y1, m1, _ := t1.Local().Date()
	y2, m2, _ := t2.Local().Date()
	return y1 == y2 && m1 == m2
}

// DaysBetweenUTC 从每日零点开始计算相差天数
func DaysBetweenUTC(before, after Time) int64 {
	t1 := ZeroTimeUTC(before)
	t2 := ZeroTimeUTC(after)
	return int64(t2.Sub(t1).Hours() / 24)
}

func DaysBetweenLocal(before, after Time) int64 {
	t1 := ZeroTimeLocal(before)
	t2 := ZeroTimeLocal(after)
	return int64(t2.Sub(t1).Hours() / 24)
}
func HoursBetween(before, after Time) int64 {
	return int64(after.Sub(before).Hours())
}
