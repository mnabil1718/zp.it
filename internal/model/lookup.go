package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/mnabil1718/zp.it/internal/cache"
)

type Lookup struct {
	ID        int       `json:"id"`
	Origin    string    `json:"origin"`
	Code      string    `json:"code"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}

type ILookup interface {
	// insert a new record to DB. Set code, origin as k, v in cache
	Insert(origin, code string) error
	// get only origin, if exists from cache. Faster for redirection
	GetOriginByCode(code string) (string, error)
	// get full row, no caching
	GetByCode(code string) (*Lookup, error)
	// increment clicks by 1 by code. No caching
	IncrementClicks(code string) error
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
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return ErrAlreadyExists
			}
		}

		return err
	}

	if err := l.cache.Set(context.Background(), code, origin, 300*time.Second); err != nil {
		return err
	}

	return nil
}

func (l *SQLiteLookup) GetByCode(code string) (*Lookup, error) {
	var lkp Lookup

	SQL := `select id, origin, code, clicks, created_at from lookup where code = ? limit 1`
	if err := l.db.QueryRow(SQL, code).Scan(
		&lkp.ID,
		&lkp.Origin,
		&lkp.Code,
		&lkp.Clicks,
		&lkp.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &lkp, nil
}

func (l *SQLiteLookup) GetOriginByCode(code string) (string, error) {
	SQL := `select origin from lookup where code = ? limit 1`
	var origin string

	v, err := l.cache.Get(context.Background(), code)
	if err == nil {
		return v, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		return "", err
	}

	if err := l.db.QueryRow(SQL, code).Scan(&origin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}

		return "", err
	}

	return origin, nil
}

func (l *SQLiteLookup) IncrementClicks(code string) error {
	SQL := `update lookup set clicks = clicks + 1 where code = ?`
	if _, err := l.db.Exec(SQL, code); err != nil {
		return err
	}

	return nil
}
