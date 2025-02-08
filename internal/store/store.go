//go:generate mockgen -source=store.go -destination=mocks/mocks.go
package store

import (
	"context"
	"github.com/Nastez/shortener/internal/app/models"
)

// Store описывает абстрактное хранилище сообщений пользователей
type Store interface {
	Bootstrap(ctx context.Context) error
	Get(ctx context.Context, id string) (string, error)
	Save(ctx context.Context, url URL) error
	SaveBatch(ctx context.Context, requestBatch models.PayloadBatch, shortURLBatch models.ResponseBodyBatch) error
}

type URL struct {
	OriginalURL string
	ShortURL    string
	GeneratedID string
}
