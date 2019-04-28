package tzdb

import (
	"testing"
	"time"
)

const (
	baseURL = "http://api.timezonedb.com"
	apiKey  = "Q34227MHXHAF"
	spath   = "v2.1/get-time-zone"
)

func Test_GetTimezone(t *testing.T) {
	now := time.Now().Unix()
	t.Logf("now: %d", now)
	client, err := NewTzdbClient(baseURL, apiKey, 3)
	if err != nil {
		t.Fatal(err)
	}

	// Test 1: Tokyo
	lon := 139.767125
	lat := 35.681236

	jisa, err := client.GetTimezone(lon, lat, now)
	if err != nil {
		t.Fatal(err)
	}
	if jisa != 32400 {
		t.Errorf("Wrong jisa, expected: %d, actual: %d", 32400, jisa)
	}

	// Avoid call limit
	time.Sleep(500 * time.Millisecond)

	// Test 2: Los Angels
	lon = -122.419416
	lat = 37.77493

	jisa, err = client.GetTimezone(lon, lat, now)
	if err != nil {
		t.Fatal(err)
	}
	if jisa != -28800 && jisa != -25200 {
		t.Errorf("Wrong jisa, expected: %d or %d, actual: %d", -28800, -25200, jisa)
	}

	time.Sleep(500 * time.Millisecond)

	// Test 3: Los Angels Winter
	pt, _ := time.LoadLocation("America/Los_Angeles")
	winterDay := time.Date(2019, 12, 1, 0, 0, 0, 0, pt).Unix()
	jisa, err = client.GetTimezone(lon, lat, winterDay)
	if err != nil {
		t.Fatal(err)
	}
	if jisa != -28800 {
		t.Errorf("Wrong jisa, expected: %d, actual: %d", -28800, jisa)
	}

	time.Sleep(500 * time.Millisecond)

	// Test 4: Los Angels Summer
	summerDay := time.Date(2019, 7, 1, 0, 0, 0, 0, pt).Unix()
	jisa, err = client.GetTimezone(lon, lat, summerDay)
	if err != nil {
		t.Fatal(err)
	}
	if jisa != -25200 {
		t.Errorf("Wrong jisa, expected: %d, actual: %d", -28800, jisa)
	}
}
