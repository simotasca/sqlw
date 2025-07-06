# sqlw

A tiny Go library for working with [`sqlx.NamedQuery`](https://pkg.go.dev/github.com/jmoiron/sqlx#DB.NamedQuery) using generics.

## Features

- ✅ Named query support
- ✅ Struct scanning
- ✅ Generics for simple usage
- ✅ Clean `Collect()` and `Scan()` helpers

## Install

```bash
go get github.com/simotasca/sqlw
```

## Usage

```go
type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type NameFilter struct {
	Name string
}

func FindByName(db *sqlx.DB, filter *NameFilter) ([]User, error) {
	return namedsqlx.CollectNamedQuery[User](
		db.NamedQuery("SELECT * FROM users WHERE name = :name", filter),
	)
}
```

