package tritondb

import (
	"database/sql"
	"fmt"

	"math/rand"

	"github.com/tada3/triton/weather/model"
)

const (
	selectByNameSql2 = "SELECT IFNULL(countryCode, ''),cityName from country_city WHERE countryName = ? OR officialName = ?"
	//selectByNameSql3      = "SELECT IFNULL(countryCode, ''),cityName from country_city WHERE countryName = ? OR officialName = ? ORDER BY RAND() LIMIT 1"
	sqlSelectCountryCity1 = "SELECT cityName from country_city WHERE countryCode = ? AND isCountry > 0"
	//sqlSelectCountryCity2      = "SELECT cityName from country_city WHERE countryCode = ? AND isCountry <> 0 ORDER BY RAND() LIMIT 1"
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
		stmtByName2, err = getDbClient().PrepareStmt(selectByNameSql2)
		if err != nil {
			log.Error("Prepare failed: %s", selectByNameSql2, err)
			return nil, false
		}
	}

	rows, err := stmtByName2.Query(cn, cn)
	if err != nil {
		stmtByName2.Close()
		stmtByName2 = nil
		log.Error("Query failed(1): %s", selectByNameSql2, err)
		return nil, false
	}
	defer rows.Close()

	cityInfos := make([]*model.CityInfo, 0)
	for rows.Next() {
		var code string
		var name string
		err := rows.Scan(&code, &name)
		if err != nil {
			stmtByName2.Close()
			stmtByName2 = nil
			log.Error("Query failed(2): %s", selectByNameSql2, err)
			return nil, false
		}
		cityInfos = append(cityInfos, &model.CityInfo{
			CountryCode: code,
			CityName:    name,
		})
	}

	len := len(cityInfos)
	if len == 0 {
		stmtByName2.Close()
		stmtByName2 = nil
		return nil, false
	}

	x := rand.Intn(len)
	return cityInfos[x], true
}

/**
func CountryName2City2(cn string) (*model.CityInfo, bool) {
	var err error
	if stmtByName2 == nil {
		stmtByName2, err = getDbClient().PrepareStmt(selectByNameSql3)
		if err != nil {
			log.Error("Prepare failed: %s", selectByNameSql3, err)
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
**/

func Country2City(code string) (*model.CityInfo, bool) {
	if stmtSelectCountryCity1 == nil {
		var pErr error
		stmtSelectCountryCity1, pErr = getDbClient().PrepareStmt(sqlSelectCountryCity1)
		if pErr != nil {
			log.Error("Prepare failed: %s", sqlSelectCountryCity1, pErr)
			return nil, false
		}
	}

	rows, err := stmtSelectCountryCity1.Query(code)
	if err != nil {
		stmtSelectCountryCity1.Close()
		stmtSelectCountryCity1 = nil
		log.Error("Query failed(1): %s", sqlSelectCountryCity1, err)
		return nil, false
	}
	defer rows.Close()

	cities := make([]string, 0)
	for rows.Next() {
		var city string
		err := rows.Scan(&city)
		if err != nil {
			stmtSelectCountryCity1.Close()
			stmtSelectCountryCity1 = nil
			log.Error("Query failed(2): %s", sqlSelectCountryCity1, err)
			return nil, false
		}
		cities = append(cities, city)
	}

	len := len(cities)
	if len == 0 {
		// Not Found
		return nil, false
	}
	x := rand.Intn(len)
	selected := cities[x]

	cityInfo := &model.CityInfo{CountryCode: code}
	cityInfo.CityName = selected
	return cityInfo, true
}

/**
func Country2City(code string) (*model.CityInfo, bool) {
	if stmtSelectCountryCity1 == nil {
		var pErr error
		stmtSelectCountryCity1, pErr = getDbClient().PrepareStmt(sqlSelectCountryCity2)
		if pErr != nil {
			log.Error("Prepare failed: %s", selectByNameSql3, pErr)
			return nil, false
		}
	}

	cityInfo := &model.CityInfo{CountryCode: code}
	var city string
	err := stmtSelectCountryCity1.QueryRow(code).Scan(&city)
	if err != nil {
		if err == sql.ErrNoRows {
			// Not Found
			return nil, false
		}
		stmtSelectCountryCity1.Close()
		stmtSelectCountryCity1 = nil
		log.Error("Query failed: %s", selectByNameSql3, err)
		return nil, false
	}
	cityInfo.CityName = city
	return cityInfo, true
}
**/

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
