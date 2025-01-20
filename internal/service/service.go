package service

import (
	"context"

	"github.com/SemenShakhray/url-shortener/internal/storage"
)

type Service struct {
	Store storage.Storer
}

type Servicer interface {
	SaveURL(ctx context.Context, urlForSave, alias string) error
	GetURL(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, alias string) error
}

func NewService(store storage.Storer) Servicer {
	return &Service{
		Store: store,
	}
}

func (s Service) SaveURL(ctx context.Context, urlForSave, alias string) error {
	return s.Store.SaveURL(ctx, urlForSave, alias)
}

func (s Service) GetURL(ctx context.Context, alias string) (string, error) {
	return s.Store.GetURL(ctx, alias)
}

func (s Service) DeleteURL(ctx context.Context, alias string) error {
	return s.Store.DeleteURL(ctx, alias)
}
