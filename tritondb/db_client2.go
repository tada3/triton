package tritondb

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tada3/triton/config"
)

const (
	dataSourceNameFmt string = "%s:%s@tcp(%s:%d)/%s"
)

var (
	defaultDbc *TritonDbClient
)

type TritonDbClient struct {
	dbType string
	db     *sql.DB
	tx     *sql.Tx
}

func init() {
	var err error
	defaultDbc, err = newTritonDbClient("sqlite3")
	if err != nil {
		panic(err)
	}
	err = defaultDbc.Open()
	if err != nil {
		panic(err)
	}
}

func getDbClient() *TritonDbClient {
	return defaultDbc
}

func (c *TritonDbClient) Open() error {
	var err error
	dsn := getDataSourceName(c.dbType)
	//log.Info("XXX Connecting2 to %s(%s)..", c.dbType, dsn)
	c.db, err = sql.Open(c.dbType, dsn)
	if err != nil {
		return err
	}
	return nil
}

// ExecTx execute sql with tx.
func (c *TritonDbClient) ExecTx(stmt string) (sql.Result, error) {
	if c.tx == nil {
		return nil, errors.New("tx does not exist")
	}
	return c.tx.Exec(stmt)
}

func (c *TritonDbClient) PrepareStmt(stmt string) (*sql.Stmt, error) {
	stmtIns, err := c.db.Prepare(stmt)
	if err != nil {
		return nil, err
	}
	return stmtIns, nil
}

func (c *TritonDbClient) PrepareStmtTx(stmt string) (*sql.Stmt, error) {
	if c.tx == nil {
		return nil, errors.New("tx does not exist")
	}
	stmtIns, err := c.tx.Prepare(stmt)
	if err != nil {
		return nil, err
	}
	return stmtIns, nil
}

func (c *TritonDbClient) Close() {
	log.Info("Closing connection to DB...")
	err := c.db.Close()
	if err != nil {
		log.Error("Failed to close DB!", err)
	}
}

func (c *TritonDbClient) BeginTx() error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	c.tx = tx
	return nil
}

func (c *TritonDbClient) CommitTx() error {
	log.Info("Committing tx..")
	if c.tx == nil {
		return errors.New("Transaction is not found.")
	}
	return c.tx.Commit()
}

func (c *TritonDbClient) RollbackTx() {
	log.Info("Rollbacking tx..")
	if c.tx == nil {
		log.Info("Transaction is not found.")
		return
	}
	err := c.tx.Rollback()
	if err != nil {
		log.Error("Rollback failed!", err)
	}
}

func newTritonDbClient(dbType string) (*TritonDbClient, error) {
	return &TritonDbClient{
		dbType: dbType,
	}, nil
}

func getDataSourceName(dbType string) string {
	cfg := config.GetConfig()
	// Assume cfg is never nil
	if dbType == "mysql" {
		return fmt.Sprintf(dataSourceNameFmt,
			cfg.MySQLUser, cfg.MySQLPasswd,
			cfg.MySQLHost, cfg.MySQLPort, cfg.MySQLDatabase)
	} else {
		return cfg.SQLiteFile
	}
}
