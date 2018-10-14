package owm

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/tada3/triton/weather/owm/db"
)

const (
	insertSql = "INSERT INTO city_list VALUES(?,?,?,?,?)"
	deleteSql = "DELETE FROM city_list"
	txSize    = 1000
)

type City struct {
	Id      int64
	Name    string
	Country string
	Coord   Coord
}

type Coord struct {
	Lon float64
	Lat float64
}

func ClearCityList() (int64, error) {
	dbc, err := db.NewOwmDbClient()
	if err != nil {
		return 0, err
	}
	err = dbc.Open()
	if err != nil {
		return 0, err
	}
	defer dbc.Close()

	return deleteCityList(dbc)
}

func LoadCityList(filepath string) (int64, error) {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return 0, err
	}
	decoder := json.NewDecoder(file)
	t, err := decoder.Token()
	if err != nil {
		return 0, err
	}
	fmt.Printf("%T: %v\n", t, t)

	dbc, err := db.NewOwmDbClient()
	if err != nil {
		return 0, err
	}
	err = dbc.Open()
	if err != nil {
		return 0, err
	}
	defer dbc.Close()

	count, err := insertCityList(decoder, dbc)
	if err != nil {
		fmt.Printf("Insert failed: %s\n", err.Error())
		fmt.Printf("  %d records were inserted.\n", count)
		return count, err
	}
	fmt.Printf("%d records were inserted.\n", count)

	t, err = decoder.Token()
	if err != nil {
		return count, err
	}

	fmt.Printf("%T: %v\n", t, t)
	return count, nil

}

func deleteCityList(dbc *db.OwmDbClient) (committed int64, err error) {
	fmt.Println("XXXXX deleteCityList 000")
	defer func() {
		fmt.Println("XXXXXXXXXXXXXX defer()")
		if r := recover(); r != nil {
			fmt.Println("Recovered in insertCityList!", r)
			dbc.RollbackTx()
			err = errors.New("unexpected error")
		}
	}()

	fmt.Println("XXXXX deleteCityList 100")

	err = dbc.BeginTx()
	if err != nil {
		return handleTxError(err, dbc, committed)
	}

	r, err := dbc.ExecTx(deleteSql)
	committed, _ = r.RowsAffected()

	err = dbc.CommitTx()
	if err != nil {
		return handleTxError(err, dbc, committed)
	}
	return committed, nil
}

func insertCityList(decoder *json.Decoder, dbc *db.OwmDbClient) (committed int64, err error) {
	fmt.Println("XXXXX insertCityList 000")
	defer func() {
		fmt.Println("XXXXXXXXXXXXXX defer()")
		if r := recover(); r != nil {
			fmt.Println("Recovered in insertCityList!", r)
			dbc.RollbackTx()
			err = errors.New("unexpected error")
		}
	}()

	fmt.Println("XXXXX insertCityList 100")

	stmt, err := startInsertTransaction(dbc)
	if err != nil {
		return handleTxError(err, dbc, committed)
	}
	count := int64(0)
	// committed := int64(0)
	for decoder.More() {
		var city City
		err := decoder.Decode(&city)
		if err != nil {
			return handleTxError(err, dbc, committed)
		}
		// fmt.Printf("city: %v\n", city.Name)

		stmt.Exec(city.Id, city.Name, city.Country, city.Coord.Lon, city.Coord.Lat)
		count++
		if count%txSize == 0 {
			fmt.Printf("Committing %d recs..\n", count)
			err = dbc.CommitTx()
			if err != nil {
				return handleTxError(err, dbc, committed)
			}
			committed += count
			count = 0

			stmt, err = startInsertTransaction(dbc)
			if err != nil {
				return handleTxError(err, dbc, committed)
			}
		}
	}
	if count > 0 {
		//dbc.RollbackTx()

		err = dbc.CommitTx()
		if err != nil {
			return handleTxError(err, dbc, committed)
		}

		committed += count
	}
	return committed, nil
}

func startInsertTransaction(dbc *db.OwmDbClient) (*sql.Stmt, error) {
	err := dbc.BeginTx()
	if err != nil {
		return nil, err
	}
	stmt, err := dbc.PrepareStmtTx(insertSql)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func handleTxError(err error, dbc *db.OwmDbClient, committed int64) (int64, error) {
	fmt.Printf("Error occurred during transaction: %s\n", err.Error())
	dbc.RollbackTx()
	return committed, err
}
