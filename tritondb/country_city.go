package tritondb

import (
	"database/sql"
	"fmt"

	"github.com/tada3/triton/weather/model"
)

const (
	selectByNameSql2           = "SELECT IFNULL(countryCode, ''),cityName from country_city WHERE countryName = ? OR officialName = ?"
	selectByNameSql3           = "SELECT IFNULL(countryCode, ''),cityName from country_city WHERE countryName = ? OR officialName = ? ORDER BY RAND() LIMIT 1"
	sqlSelectCountryCity1      = "SELECT cityName from country_city WHERE countryCode = ? AND isCountry > 0"
	sqlSelectCountryCity2      = "SELECT cityName from country_city WHERE countryCode = ? AND isCountry <> 0 ORDER BY RAND() LIMIT 1"
	selectCountryNameByCodeSQL = "SELECT countryName from country_city WHERE countryCode = ? AND isCountry > 0"
)

var (
	stmtByName             *sql.Stmt
	stmtByCode             *sql.Stmt
	stmtCountryNameByCode  *sql.Stmt
	stmtByName2            *sql.Stmt
	stmtSelectCountryCity1 *sql.Stmt
)

func CountryName2City2(cn string) (*model.CityInfo, bool) {
	var err error
	if stmtByName2 == nil {
		stmtByName2, err = getDbClient().PrepareStmt(selectByNameSql3)
		if err != nil {
			fmt.Printf("ERROR! Prepare failed: %s, stmt=%v\n", err.Error(), selectByNameSql3)
			return nil, false
		}
	}

	cityInfo := &model.CityInfo{}
	err = stmtByName2.QueryRow(cn, cn).Scan(&(cityInfo.CountryCode),
		&(cityInfo.CityName))
	if err != nil {
		if err != sql.ErrNoRows {
			stmtByName2.Close()
			stmtByName2 = nil
		}
		return nil, false
	}
	return cityInfo, true
}

func Country2City(code string) (*model.CityInfo, bool, error) {
	if stmtSelectCountryCity1 == nil {
		var pErr error
		stmtSelectCountryCity1, pErr = getDbClient().PrepareStmt(sqlSelectCountryCity2)
		if pErr != nil {
			return nil, false, pErr
		}
	}

	cityInfo := &model.CityInfo{CountryCode: code}
	var city string
	err := stmtSelectCountryCity1.QueryRow(code).Scan(&city)
	if err != nil {
		if err == sql.ErrNoRows {
			// Not Found
			return nil, false, nil
		}
		stmtSelectCountryCity1.Close()
		stmtSelectCountryCity1 = nil
		return nil, false, err
	}
	cityInfo.CityName = city
	return cityInfo, true, nil
}

func CountryCode2CountryName(code string) (string, bool) {
	if stmtCountryNameByCode == nil {
		var pErr error
		stmtCountryNameByCode, pErr = getDbClient().PrepareStmt(selectCountryNameByCodeSQL)
		if pErr != nil {
			// ERROR!
			fmt.Printf("ERROR! DB ERROR: %s\n", pErr.Error())
			return "", false
		}
	}

	var country string
	err := stmtCountryNameByCode.QueryRow(code).Scan(&country)
	if err != nil {
		if err != sql.ErrNoRows {
			// ERROR!
			fmt.Printf("ERROR! Query failed: %s\n", err.Error())
			stmtCountryNameByCode.Close()
			stmtCountryNameByCode = nil
		}
		return "", false
	}
	return country, true
}
