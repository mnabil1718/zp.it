package model

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
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
	Insert(origin, code string) error
	GetByCode(code string) (*Lookup, error)
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
	var lkp Lookup
	SQL := `insert into lookup (origin, code) values (?, ?)
			returning id, origin, code, clicks, created_at`

	if err := l.db.QueryRow(SQL, origin, code).Scan(
		&lkp.ID,
		&lkp.Origin,
		&lkp.Code,
		&lkp.Clicks,
		&lkp.CreatedAt,
	); err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return ErrAlreadyExists
			}
		}

		return err
	}

	if err := l.cache.Set(context.Background(), code, lkp, 300*time.Second); err != nil {
		return err
	}

	return nil
}

func (l *SQLiteLookup) GetByCode(code string) (*Lookup, error) {
	var lkp Lookup

	err := l.cache.Get(context.Background(), code, &lkp)
	if err == nil {
		slog.Info("cache hit", "url", lkp.Origin, "code", lkp.Code, "clicks", lkp.Clicks, "created_at", lkp.CreatedAt)
		return &lkp, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		return nil, err
	}

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

	// save back to cache
	if err = l.cache.Set(context.Background(), code, lkp, 300*time.Second); err != nil {
		return nil, err
	}

	return &lkp, nil
}
