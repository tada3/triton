package tritondb

import (
	"database/sql"
)

const (
	selectByNameSql = "SELECT cityName from country_city WHERE countryName = ? OR officialName = ?"
	selectByCodeSql = "SELECT cityName from country_city WHERE countryCode = ?"
)

var (
	stmtByName *sql.Stmt
	stmtByCode *sql.Stmt
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
