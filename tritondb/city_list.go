package tritondb

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	selectPreferredCitySQL = "SELECT id from preferred_city WHERE name = ? ORDER BY priority DESC"
	selectCityListSQL      = "SELECT id from city_list WHERE name = ?"
	removeShiSQL           = "UPDATE city_list SET name = TRIM(TRAILING '-shi' FROM name) WHERE country = 'JP' AND name LIKE '%-shi'"
)

var (
	stmtP *sql.Stmt
	stmtC *sql.Stmt
)

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
