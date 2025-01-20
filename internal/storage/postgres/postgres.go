package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/SemenShakhray/url-shortener/internal/config"
	"github.com/SemenShakhray/url-shortener/internal/storage"

	"github.com/lib/pq"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) storage.Storer {
	return &Store{
		DB: db,
	}
}

func Connect(cfg config.Config) (*sql.DB, error) {
	op := "storage.postgres.Connect"
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Pass, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

func (s *Store) SaveURL(ctx context.Context, urlForSave, alias string) error {
	op := "storage.postgrese.SaveURL"

	stmt, err := s.DB.Prepare("INSERT INTO url (url, alias, create_at) VALUES ($1, $2, CURRENT_TIMESTAMP)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, urlForSave, alias)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrURLOrAliasExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if n == 0 {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Store) GetURL(ctx context.Context, alias string) (string, error) {
	op := "storage.postgres.GetURL"

	stmt, err := s.DB.Prepare("SELECT url FROM url WHERE alias=$1")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var url string

	err = stmt.QueryRowContext(ctx, alias).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Store) DeleteURL(ctx context.Context, alias string) error {
	op := "storage.postgres.DeleteURL"

	stmt, err := s.DB.Prepare("DELETE FROM url WHERE alias=$1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if n == 0 {
		return fmt.Errorf("%s: failed delete url", op)
	}
	return nil
}
