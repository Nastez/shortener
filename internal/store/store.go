//go:generate mockgen -source=store.go -destination=mocks/mocks.go
package store

import (
	"context"
)

// Store описывает абстрактное хранилище сообщений пользователей
type Store interface {
	Bootstrap(ctx context.Context) error
	Get(ctx context.Context, id string) (string, error)
	Save(ctx context.Context, url URL) error
}

type URL struct {
	OriginalURL string
	ShortURL    string
	GeneratedID string
}
