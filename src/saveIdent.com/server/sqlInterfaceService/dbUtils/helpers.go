package dbUtils

import (
	"database/sql"
	"errors"
)

func ExecuteLineCheckQuery(query string, callingFunction string, db *sql.DB) error{
	result, err := db.Exec(query)
	if err != nil{
		return err
	}

	if rows, _ := result.RowsAffected(); rows != 1{
		err = errors.New("no lines affected")
		return err
	}

	return nil
}