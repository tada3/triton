package owm

import (
	"database/sql"

	"github.com/tada3/triton/weather/owm/db"
)

const (
	selectSql = "SELECT id from city_list WHERE name = ?"
)

var dbc *db.OwmDbClient
var stmt *sql.Stmt

func init() {
	var err error
	dbc, err = db.NewOwmDbClient()
	if err != nil {
		panic(err)
	}
}

func GetCityID(city string) (int64, bool, error) {
	if stmt == nil {
		var pErr error
		stmt, pErr = dbc.PrepareStmt(selectSql)
		if pErr != nil {
			return -1, false, pErr
		}
	}

	var id int64
	err := stmt.QueryRow(city).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Not Found
			return -1, false, nil
		}
		stmt.Close()
		stmt = nil
		return -1, false, err
	}

	return id, true, nil
}
