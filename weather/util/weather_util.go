package util

import (
	"fmt"
	"math"
)

const (
	tempStrFormatP  string = "%d度"
	tempStrFormatN  string = "氷点下%d度"
	tempRangeFormat string = "%sから%s"
)

func MarumeTemp(t float64) int64 {
	if t < 0 {
		// '通常は地上1.25～2.0mの大気の温度を摂氏（℃）単位で表す。度の単位に丸めるときは十分位を四捨五入するが、０度未満は五捨六入する。'
		// by 気象庁
		return int64(math.Ceil(t - 0.5))
	}
	return int64(math.Floor(t + 0.5))
}

func GetTempStr(t int64) string {
	if t < 0 {
		return fmt.Sprintf(tempStrFormatN, -1*t)
	}
	return fmt.Sprintf(tempStrFormatP, t)
}

func GetTempRangeStr(tmin, tmax int64) string {
	if tmin == tmax {
		return GetTempStr(tmin)
	}
	return fmt.Sprintf(tempRangeFormat, GetTempStr(tmin), GetTempStr(tmax))
}
