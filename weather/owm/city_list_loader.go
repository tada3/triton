package owm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/tada3/triton/weather/owm/db"
)

const (
	TX_SIZE = 100
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

func insertCityList(decoder *json.Decoder, dbc *db.OwmDbClient) (commited int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in insertCityList!", r)
			dbc.RollbackTx()
			err = errors.New("Unexpected error!")
		}
	}()

	stmt, err := dbc.PrepareStmt("INSERT INTO city_list VALUES(?,?,?,?)")
	if err != nil {
		return commited, err
	}

	err = dbc.BeginTx()
	count := int64(0)
	committed := int64(0)
	for decoder.More() {
		fmt.Println("XXXXXXX")
		var city City

		err := decoder.Decode(&city)
		if err != nil {
			return committed, err
		}
		fmt.Printf("city: %v\n", city.Name)

		stmt.Exec(city.Id, city.Name, city.Country, city.Coord.Lon, city.Coord.Lat)
		count++
		if count%TX_SIZE == 0 {
			err = dbc.CommitTx()
			if err != nil {
				return committed, err
			}
			count = 0
			committed += count
		}
	}
	if count > 0 {
		err = dbc.CommitTx()
		if err != nil {
			return committed, err
		}
		committed += count
	}
	return committed, nil
}
