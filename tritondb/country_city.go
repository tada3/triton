package tritondb

import (
	"database/sql"
)

const (
	selectByNameSql = "SELECT cityName from country_city WHERE countryName = ? OR officialName = ?"
	selectByCodeSql = "SELECT cityName from country_city WHERE countryCode = ?"
)

var (
	stmt *sql.Stmt
)

func CountryName2City(cn string) (string, bool, error) {
	if stmt == nil {
		var pErr error
		stmt, pErr = getDbClient().PrepareStmt(selectByNameSql)
		if pErr != nil {
			return "", false, pErr
		}
	}

	var city string
	err := stmt.QueryRow(cn).Scan(&city)
	if err != nil {
		if err == sql.ErrNoRows {
			// Not Found
			return "", false, nil
		}
		stmt.Close()
		stmt = nil
		return "", false, err
	}

	return city, true, nil
}

func CountryCode2City(code string) (string, bool, error) {
	if stmt == nil {
		var pErr error
		stmt, pErr = getDbClient().PrepareStmt(selectByCodeSql)
		if pErr != nil {
			return "", false, pErr
		}
	}

	var city string
	err := stmt.QueryRow(code).Scan(&city)
	if err != nil {
		if err == sql.ErrNoRows {
			// Not Found
			return "", false, nil
		}
		stmt.Close()
		stmt = nil
		return "", false, err
	}

	return city, true, nil
}
