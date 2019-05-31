package util

import (
	"fmt"
	"time"
)

const (
	dayFormat = "%d日"
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

func GetTomorrowNoonUt(now int64, offset int64) (int64, *time.Time) {
	localTime := ut2Time(now, offset, 0)

	// 0 -> +12, 1 -> +11 + 24, ..., 12 -> +24, ..., 23 -> +13
	var delta time.Duration = 12
	if localTime.Hour() > 0 {
		delta += time.Duration(24 - localTime.Hour())
	}

	tomorrowNoon := localTime.Add(time.Hour * delta)
	rounded := tomorrowNoon.Round(time.Hour)
	return time2Ut(rounded, offset, 0), &rounded
}

func GetDayStr(d int) string {
	if d == 1 {
		return "ついたち"
	}
	switch d {
	case 1:
		return "ついたち"
	case 2:
		return "ふつか"
	case 3:
		return "みっか"
	case 4:
		return "よっか"
	case 5:
		return "いつか"
	case 6:
		return "むいか"
	case 7:
		return "なのか"
	case 8:
		return "ようか"
	case 9:
		return "ここのか"
	case 10:
		return "とおか"
	default:
		return fmt.Sprintf(dayFormat, d)
	}
}
