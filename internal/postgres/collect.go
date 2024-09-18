package postgres

import (
	"github.com/jackc/pgx/v5"
)

// AppendRowsNonNil iterates through rows, calling fn for each row, and appending the results into a slice of T if the result is not nil.
//
// This function closes the rows automatically on return.
func AppendRowsNonNil[T comparable, S ~[]T](slice S, rows pgx.Rows, fn pgx.RowToFunc[T]) (S, error) {
	defer rows.Close()

	for rows.Next() {
		value, err := fn(rows)
		if err != nil {
			return nil, err
		}
		var t T
		if value == t {
			continue
		}
		slice = append(slice, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return slice, nil
}

// CollectRowsNonNil iterates through rows, calling fn for each row, and collecting the results into a slice of T if the result is not nil.
//
// This function closes the rows automatically on return.
func CollectRowsNonNil[T comparable](rows pgx.Rows, fn pgx.RowToFunc[T]) ([]T, error) {
	return AppendRowsNonNil([]T{}, rows, fn)
}
