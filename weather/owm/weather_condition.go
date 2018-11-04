package owm

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	csvFilePath = "weather/owm/csv/weather_condition.csv"
	unknown     = "不明"
)

var (
	wcMap = map[int64]string{}
)

func init() {
	f, err := os.Open(csvFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = loadCsvFile(f)
	if err != nil {
		panic(err)
	}
}

func GetWeatherCondition(code int64) string {
	wc, ok := wcMap[code]
	if !ok {
		fmt.Printf("Unknown code: %d\n", code)
		return unknown
	}
	return wc
}

func loadCsvFile(f *os.File) error {
	r := csv.NewReader(bufio.NewReader(f))
	r.Comment = '#'
	r.TrimLeadingSpace = true
	count := 0
	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if len(l) != 4 {
			return fmt.Errorf("Invalid format: %v", l)
		}
		id, err := strconv.ParseInt(l[0], 10, 64)
		if err != nil {
			return err
		}
		wcMap[id] = l[2]
		count++
	}
	fmt.Printf("Loaded %d records.\n", count)
	return nil
}
