package tritondb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tada3/triton/weather/model"
)

const (
	selectPreferredCitySQL          = "SELECT id from preferred_city WHERE name = ? ORDER BY priority DESC"
	selectCityListSQL               = "SELECT id from city_list WHERE name = ?"
	removeShiSQL                    = "UPDATE city_list SET name = TRIM(TRAILING '-shi' FROM name) WHERE country = 'JP' AND name LIKE '%-shi'"
	selectPreferredCitySQL2         = "SELECT id from preferred_city WHERE name=? AND country=? ORDER BY priority DESC"
	selectPreferredCityNoCountrySQL = "SELECT id, country from preferred_city WHERE name=? ORDER BY priority DESC"
	selectCityListSQL2              = "SELECT id from city_list WHERE name=? AND country=?"
	selectCityListNoCountrySQL      = "SELECT id, country from city_list WHERE name=?"
)

var (
	stmtP    *sql.Stmt
	stmtC    *sql.Stmt
	stmtP2   *sql.Stmt
	stmtP2NC *sql.Stmt
	stmtC2   *sql.Stmt
	stmtC2NC *sql.Stmt
)

func getCityIDFromPreferredCity(cityName, countryCode string) (int64, bool) {
	var err error
	if stmtP2 == nil {
		stmtP2, err = getDbClient().PrepareStmt(selectPreferredCitySQL2)
		if err != nil {
			fmt.Printf("ERROR! Prepare failed: %s, stmt=%v\n", err.Error(), selectPreferredCitySQL2)
			return -1, false
		}
	}

	var id int64
	err = stmtP2.QueryRow(cityName, countryCode).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			// Error
			fmt.Printf("ERROR! Query failed: %s, stmt=%v\n", err.Error(), stmtP2)
			stmtP2.Close()
			stmtP2 = nil
		}
		return 0, false
	}
	return id, true
}

func getCityIDFromPreferredCityNoCountry(cityName string) (int64, string, bool) {
	var err error
	if stmtP2NC == nil {
		stmtP2NC, err = getDbClient().PrepareStmt(selectPreferredCityNoCountrySQL)
		if err != nil {
			fmt.Printf("ERROR! Prepare failed: %s, stmt=%v\n", err.Error(), selectPreferredCityNoCountrySQL)
			return 0, "", false
		}
	}

	var id int64
	var code string
	err = stmtP2NC.QueryRow(cityName).Scan(&id, &code)
	if err != nil {
		if err != sql.ErrNoRows {
			// Error
			fmt.Printf("ERROR! Query failed: %s, stmt=%v\n", err.Error(), stmtP2NC)
			stmtP2NC.Close()
			stmtP = nil
		}
		return 0, "", false
	}
	return id, code, true
}

func getCityIDFromCityList(cityName, countryCode string) (int64, bool) {
	var err error
	if stmtC2 == nil {
		stmtC2, err = getDbClient().PrepareStmt(selectCityListSQL2)
		if err != nil {
			fmt.Printf("ERROR! Prepare failed: %s, stmt=%v\n", err.Error(), selectCityListSQL2)
			return 0, false
		}
	}

	var id int64
	err = stmtC2.QueryRow(cityName, countryCode).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			// Error
			fmt.Printf("ERROR! Query failed: %s, stmt=%v\n", err.Error(), stmtC2)
			stmtC2.Close()
			stmtC2 = nil
		}
		return 0, false
	}
	return id, true
}

func getCityIDFromCityListNoCountry(cityName string) (int64, string, bool) {
	var err error
	if stmtC2NC == nil {
		stmtC2NC, err = getDbClient().PrepareStmt(selectCityListNoCountrySQL)
		if err != nil {
			fmt.Printf("ERROR! Prepare failed: %s, stmt=%v\n", err.Error(), selectCityListNoCountrySQL)
			return 0, "", false
		}
	}

	var id int64
	var code string
	err = stmtC2NC.QueryRow(cityName).Scan(&id, &code)
	if err != nil {
		if err != sql.ErrNoRows {
			// Error
			fmt.Printf("ERROR! Query failed: %s, stmt=%v\n", err.Error(), stmtC2NC)
			stmtC2NC.Close()
			stmtC2NC = nil
		}
		return 0, "", false
	}
	return id, code, true

}

// GetCityID2 get city ID for the specified city name from DB.
func GetCityID2(city *model.CityInfo) (int64, string, bool) {

	if city.CountryCode != "" {
		// By cityName and countryCode
		id, found := getCityIDFromPreferredCity(city.CityNameEN, city.CountryCode)
		if found {
			return id, city.CountryCode, true
		}
		id, found = getCityIDFromCityList(city.CityNameEN, city.CountryCode)
		if found {
			return id, city.CountryCode, true
		}
	} else {
		// By cityName only
		id, code, found := getCityIDFromPreferredCityNoCountry(city.CityNameEN)
		if found {
			return id, code, true
		}
		id, code, found = getCityIDFromCityListNoCountry(city.CityNameEN)
		if found {
			return id, code, true
		}
	}
	return 0, "", false
}

// GetCityID get city ID for the specified city name from DB.
func GetCityID(city string) (int64, bool, error) {
	// 1. Check preferred_city
	var err error
	if stmtP == nil {
		stmtP, err = getDbClient().PrepareStmt(selectPreferredCitySQL)
		if err != nil {
			return -1, false, err
		}
	}

	var id int64
	err = stmtP.QueryRow(city).Scan(&id)
	if err == nil {
		return id, true, nil
	}

	if err != sql.ErrNoRows {
		// Error
		stmtP.Close()
		stmtP = nil
	}

	// 2. Get id from city_list
	if stmtC == nil {
		stmtC, err = getDbClient().PrepareStmt(selectCityListSQL)
		if err != nil {
			return -1, false, err
		}
	}

	err = stmtC.QueryRow(city).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Not Found
			return -1, false, nil
		}
		stmtC.Close()
		stmtC = nil
		return -1, false, err
	}

	return id, true, nil
}

// RemoveShiFromJPCities removes '-shi' from names of JP cities.
// It is easier to match city names without it.
func RemoveShiFromJPCities() (count int64, err error) {
	dbc := getDbClient()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RemoveShiFromJPCities()!", r)
			dbc.RollbackTx()
			count = 0
			err = errors.New("unexpected error")
		}
	}()

	err = dbc.BeginTx()
	if err != nil {
		return handleTxError(err, dbc)
	}

	result, err := dbc.ExecTx(removeShiSQL)
	if err != nil {
		return handleTxError(err, dbc)
	}
	count, _ = result.RowsAffected()

	err = dbc.CommitTx()
	if err != nil {
		return handleTxError(err, dbc)
	}
	return count, nil
}

func handleTxError(err error, dbc *TritonDbClient) (int64, error) {
	fmt.Printf("Error occurred during transaction: %s\n", err.Error())
	dbc.RollbackTx()
	return 0, err
}
