package utils

import (
	"database/sql"
)

// CheckRowsAffected verifies that an Exec result affected at least one row.
// Returns ErrNotFound if 0 rows were affected.
func CheckRowsAffected(res sql.Result) error {
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}
