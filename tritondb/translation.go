package tritondb

import (
	"database/sql"
)

const (
	sqlSelectTranslation = "SELECT dst from translation WHERE src = ?"
)

var (
	stmtSelectTranslation *sql.Stmt
)

func TranslateByDB(w string) (string, bool, error) {
	if stmtSelectTranslation == nil {
		var pErr error
		stmtSelectTranslation, pErr = getDbClient().PrepareStmt(sqlSelectTranslation)
		if pErr != nil {
			return "", false, pErr
		}
	}

	var dst string
	err := stmtSelectTranslation.QueryRow(w).Scan(&dst)
	if err != nil {
		if err == sql.ErrNoRows {
			// Not Found
			return "", false, nil
		}
		stmtSelectTranslation.Close()
		stmtSelectTranslation = nil
		return "", false, err
	}
	return dst, true, nil
}
