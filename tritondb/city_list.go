package tritondb

import (
	"errors"
	"fmt"
)

const (
	removeShiSQL = "UPDATE city_list SET name = TRIM(TRAILING '-shi' FROM name) WHERE country = 'JP' AND name LIKE '%-shi'"
)

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
