package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tada3/triton/config"
)

const (
	dataSourceNameFmt string = "%s:%s@tcp(%s:%d)/%s"
)

var (
	defaultDbc *OwmDbClient
)

type OwmDbClient struct {
	db *sql.DB
	tx *sql.Tx
}

func init() {
	var err error
	defaultDbc, err = NewOwmDbClient()
	if err != nil {
		panic(err)
	}
	err = defaultDbc.Open()
	if err != nil {
		panic(err)
	}
}

func GetDbClient() *OwmDbClient {
	return defaultDbc
}

func NewOwmDbClient() (*OwmDbClient, error) {
	return &OwmDbClient{}, nil
}

func (c *OwmDbClient) Open() error {
	var err error
	dsn := getDataSourceName()
	fmt.Printf("INFO Connecting to MySQL(%s)..\n", dsn)
	c.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return nil
}

// ExecTx execute sql with tx.
func (c *OwmDbClient) ExecTx(stmt string) (sql.Result, error) {
	if c.tx == nil {
		return nil, errors.New("tx does not exist")
	}
	return c.tx.Exec(stmt)
}

func (c *OwmDbClient) PrepareStmt(stmt string) (*sql.Stmt, error) {
	stmtIns, err := c.db.Prepare(stmt)
	if err != nil {
		return nil, err
	}
	return stmtIns, nil
}

func (c *OwmDbClient) PrepareStmtTx(stmt string) (*sql.Stmt, error) {
	if c.tx == nil {
		return nil, errors.New("tx does not exist")
	}
	stmtIns, err := c.tx.Prepare(stmt)
	if err != nil {
		return nil, err
	}
	return stmtIns, nil
}

func (c *OwmDbClient) Close() {
	err := c.db.Close()
	if err != nil {
		fmt.Printf("Failed to close DB: %s\n", err.Error())
	}
}

func (c *OwmDbClient) BeginTx() error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	c.tx = tx
	return nil
}

func (c *OwmDbClient) CommitTx() error {
	fmt.Println("Committing tx..")
	if c.tx == nil {
		return errors.New("Transaction is not found.")
	}
	return c.tx.Commit()
}

func (c *OwmDbClient) RollbackTx() {
	fmt.Println("Rollbacking tx..")
	if c.tx == nil {
		fmt.Println("Transaction is not found.")
		return
	}
	err := c.tx.Rollback()
	if err != nil {
		fmt.Printf("Rollback failed: %s\n", err.Error())
	}
}

func getDataSourceName() string {
	cfg := config.GetConfig()
	// Assume cfg is never nil
	return fmt.Sprintf(dataSourceNameFmt,
		cfg.MySQLUser, cfg.MySQLPasswd,
		cfg.MySQLHost, cfg.MySQLPort, cfg.MySQLDatabase)
}
