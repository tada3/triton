package model

type CurrentWeather struct {
	Weather     string
	Temp        int64
	CountryCode string
}

type CityInfo struct {
	CountryCode string
	CityName    string
	CityNameEN  string
	CityID      int64
}
