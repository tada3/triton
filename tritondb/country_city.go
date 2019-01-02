package tritondb

import (
	"database/sql"
	"fmt"

	"github.com/tada3/triton/weather/model"
)

const (
	selectByNameSql            = "SELECT cityName from country_city WHERE countryName = ? OR officialName = ?"
	selectByNameSql2           = "SELECT IFNULL(countryCode, ''),cityName from country_city WHERE countryName = ? OR officialName = ?"
	selectByCodeSql            = "SELECT cityName from country_city WHERE countryCode = ?"
	selectCountryNameByCodeSQL = "SELECT countryName from country_city WHERE countryCode = ? AND isCountry > 0"
)

var (
	stmtByName            *sql.Stmt
	stmtByCode            *sql.Stmt
	stmtCountryNameByCode *sql.Stmt
	stmtByName2           *sql.Stmt
)

func CountryName2City(cn string) (string, bool, error) {
	if stmtByName == nil {
		var pErr error
		stmtByName, pErr = getDbClient().PrepareStmt(selectByNameSql)
		if pErr != nil {
			return "", false, pErr
		}
	}

	var city string
	err := stmtByName.QueryRow(cn, cn).Scan(&city)
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

func CountryName2City2(cn string) (*model.CityInfo, bool) {
	var err error
	if stmtByName2 == nil {
		stmtByName2, err = getDbClient().PrepareStmt(selectByNameSql2)
		if err != nil {
			fmt.Printf("ERROR! Prepare failed: %s, stmt=%v\n", err.Error(), selectByNameSql2)
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

func CountryCode2City(code string) (string, bool, error) {
	if stmtByCode == nil {
		var pErr error
		stmtByCode, pErr = getDbClient().PrepareStmt(selectByCodeSql)
		if pErr != nil {
			return "", false, pErr
		}
	}

	var city string
	err := stmtByCode.QueryRow(code).Scan(&city)
	if err != nil {
		if err == sql.ErrNoRows {
			// Not Found
			return "", false, nil
		}
		stmtByCode.Close()
		stmtByCode = nil
		return "", false, err
	}

	return city, true, nil
}

func CountryCode2City2(code string) (*model.CityInfo, bool, error) {
	if stmtByCode == nil {
		var pErr error
		stmtByCode, pErr = getDbClient().PrepareStmt(selectByCodeSql)
		if pErr != nil {
			return nil, false, pErr
		}
	}

	cityInfo := &model.CityInfo{CountryCode: code}
	var city string
	err := stmtByCode.QueryRow(code).Scan(&city)
	if err != nil {
		if err == sql.ErrNoRows {
			// Not Found
			return nil, false, nil
		}
		stmtByCode.Close()
		stmtByCode = nil
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
