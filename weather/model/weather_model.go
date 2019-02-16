package model

type CurrentWeather struct {
	Weather     string
	Temp        int64
	TempStr     string
	CountryCode string
}

type CityInfo struct {
	CountryCode string
	CityName    string
	CityNameEN  string
	CityID      int64
}

func (a *CityInfo) Clone() *CityInfo {
	copy := *a
	return &copy
}
