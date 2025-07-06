package sqlw

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// wraps a named query and returns a single result of type T
func NamedQueryOne[T any](rows *sqlx.Rows, err error) (*T, error) {
	if err != nil {
		return nil, fmt.Errorf("exec named query one: %w", err);
	}
	if !rows.Next() {
        return nil, sql.ErrNoRows
    }
    var dest = new(T)
    err = rows.Scan(dest)
    if err != nil {
        return nil, fmt.Errorf("scan named query one: %w", err)
    }
    return dest, nil
}

// wraps a named query and returns a SqlxMapper[T]
func NamedQueryMapper[T any](rows *sqlx.Rows, err error) (*SqlxMapper[T], error) {
	if err != nil {
		return nil, err;
	}
	return &SqlxMapper[T]{rows}, nil
}

// wraps a named query and collects all results into a slice of T
func CollectNamedQuery[T any](rows *sqlx.Rows, err error) ([]T, error) {
	mapper, err := NamedQueryMapper[T](rows, err) 
	if err != nil {
		return nil, err;
	}
	return mapper.Collect()
}

// wraps sqlx.Rows to provide generic iter-like mapper over type T
type SqlxMapper[T any] struct {
	*sqlx.Rows
}

// returns true if there is another row to read
func (m *SqlxMapper[T]) Next() bool {
	return m.Rows.Next()
}

// reads the current row into a new instance of T
func (m *SqlxMapper[T]) Scan() (*T, error) {
	dest := new(T)
	err := m.Rows.StructScan(dest)
	return dest, err
}

// closes the underlying rows
func (m *SqlxMapper[T]) Close() error {
	return m.Rows.Close()
}

// reads all remaining rows and returns them as a slice of T
func (m *SqlxMapper[T]) Collect() ([]T, error) {
	defer m.Close()
	result := make([]T, 0)
	for m.Next() {
		row, err := m.Scan()
		if err != nil {
			return nil, fmt.Errorf("mapper scan: %w", err)
		}
		result = append(result, *row)
	}
	return result, nil
}