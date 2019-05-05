package timezone

import (
	"fmt"
	"time"

	"github.com/tada3/triton/timezone/tzdb"
	"github.com/tada3/triton/weather/util"
)

const (
	baseURL = "http://api.timezonedb.com"
	apiKey  = "Q34227MHXHAF"
	spath   = "v2.1/get-time-zone"
)

var client *tzdb.TzdbClient

func init() {
	var err error
	client, err = tzdb.NewTzdbClient(baseURL, apiKey, 3)
	if err != nil {
		panic(err)
	}

}

// GetTomorrowNoon returns the timestamp of tomorrow noon at the specified
// place.
// To be exactly, you need to calculate the case where the border of
// DST/Non-DST comes between now and tomorrow noon. But I skipped it.
func GetTomorrowNoon(lon, lat float64, now int64) (int64, *time.Time) {
	jisa, err := client.GetTimezone(lon, lat, now)
	if err != nil {
		fmt.Printf("ERROR: Failed to get timezone(%f, %f, %d): %s Use JST (32400).", lon, lat, now, err.Error())
		// Use JST
		jisa = 32400
	}
	return util.GetTomorrowNoonUt(now, jisa)
}
