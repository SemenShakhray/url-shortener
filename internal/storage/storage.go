package storage

import (
	"context"
	"errors"
)

var (
	ErrURLOrAliasExists = errors.New("url or alias already exists")
	ErrURLNotFound      = errors.New("url not founded")
)

type Storer interface {
	SaveURL(ctx context.Context, urlForSave, alias string) error
	GetURL(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, alias string) error
}
