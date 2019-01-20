package tritondb

import (
	"database/sql"
	"fmt"

	"github.com/tada3/triton/weather/model"
)

const (
	sqlSelectPoiCity1 = "SELECT cityName from poi_city WHERE (name = ? OR name2 = ?) AND countryCode = ?"
	sqlSelectPoiCity2 = "SELECT cityName, countryCode from poi_city WHERE (name = ? OR name2 = ?)"
)

var (
	stmtSelectPoiCity1 *sql.Stmt
	stmtSelectPoiCity2 *sql.Stmt
)

// Poi2City gets a City for the specified POI.
func Poi2City(poi string, cityInfo *model.CityInfo) (*model.CityInfo, bool, error) {
	var cityName, countryCode string
	var found bool
	var err error
	if cityInfo == nil {
		cityInfo = &model.CityInfo{}
	}
	if cityInfo.CountryCode != "" {
		cityName, found, err = getPoiCityByNameAndCountry(poi, cityInfo.CountryCode)
		if err != nil {
			return cityInfo, false, err
		}
		if !found {
			return cityInfo, false, nil
		}

		cityInfo.CityName = cityName
	} else {
		cityName, countryCode, found, err = getPoiCityByName(poi)
		if err != nil {
			return cityInfo, false, err
		}
		if !found {
			return cityInfo, false, nil
		}
		cityInfo.CityName = cityName
		cityInfo.CountryCode = countryCode
	}
	return cityInfo, true, nil
}

func getPoiCityByNameAndCountry(poi string, cc string) (string, bool, error) {
	if stmtSelectPoiCity1 == nil {
		var pErr error
		stmtSelectPoiCity1, pErr = getDbClient().PrepareStmt(sqlSelectPoiCity1)
		if pErr != nil {
			return "", false, pErr
		}
	}

	var city string
	err := stmtSelectPoiCity1.QueryRow(poi, poi, cc).Scan(&city)
	if err != nil {
		if err == sql.ErrNoRows {
			// Not Found
			return "", false, nil
		}
		stmtByName.Close()
		stmtByName = nil
		return "", false, err
	}

	return city, true, nil
}

func getPoiCityByName(poi string) (string, string, bool, error) {
	fmt.Println("XXX getPoiCityByName 000")
	if stmtSelectPoiCity2 == nil {
		var pErr error
		stmtSelectPoiCity2, pErr = getDbClient().PrepareStmt(sqlSelectPoiCity2)
		if pErr != nil {
			return "", "", false, pErr
		}
	}

	var city, country string
	err := stmtSelectPoiCity2.QueryRow(poi, poi).Scan(&city, &country)
	if err != nil {
		fmt.Printf("XXX stmt = %v\n", stmtSelectPoiCity2)
		if err == sql.ErrNoRows {
			// Not Found
			return "", "", false, nil
		}
		stmtSelectPoiCity2.Close()
		stmtSelectPoiCity2 = nil
		return "", "", false, err
	}
	return city, country, true, nil
}
