package util

import (
	"time"
)

func ut2Time(t int64, offset int64, dst int64) time.Time {
	t1 := t + offset + dst
	return time.Unix(t1, 0).UTC()
}

func time2Ut(t time.Time, offset int64, dst int64) int64 {
	return t.Unix() - offset - dst
}

func getTomorrowRange(now int64, offset int64, dst int64) (int64, int64) {
	nowT := ut2Time(now, offset, dst)
	year, month, day := nowT.Date()
	tomorrowT := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Add(time.Duration(24) * time.Hour)
	start := time2Ut(tomorrowT, offset, dst)
	end := start + (24 * 60 * 60)
	return start, end
}
