package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/mnabil1718/zp.it/internal/cache"
)

type Lookup struct {
	id     int
	origin string
	code   string
}

type ILookup interface {
	Insert(origin, code string) error
	GetByCode(code string) (string, error)
}

type SQLiteLookup struct {
	db    *sql.DB
	cache cache.ICache
}

func NewSQliteLookup(db *sql.DB, cache cache.ICache) *SQLiteLookup {
	return &SQLiteLookup{
		db:    db,
		cache: cache,
	}
}

func (l *SQLiteLookup) Insert(origin, code string) error {
	SQL := `insert into lookup (origin, code) values (?, ?)`
	if _, err := l.db.Exec(SQL, origin, code); err != nil {
		return err
	}

	if err := l.cache.Set(context.Background(), code, origin, 300*time.Second); err != nil {
		return err
	}

	return nil
}

func (l *SQLiteLookup) GetByCode(code string) (string, error) {
	SQL := `select origin from lookup where code = ? limit 1`
	var origin string

	res, err := l.cache.Get(context.Background(), code)
	if err == nil {
		return res, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		return "", err
	}

	if err := l.db.QueryRow(SQL, code).Scan(&origin); err != nil {
		return "", err
	}

	return origin, nil
}
