//go:generate mockgen -source=store.go -destination=mocks/mocks.go
package store

import (
	"context"
	"errors"
	"github.com/Nastez/shortener/internal/app/models"
)

// ErrConflict указывает на конфликт данных в хранилище.
var ErrConflict = errors.New("data conflict")

// Store описывает абстрактное хранилище сообщений пользователей
type Store interface {
	Get(ctx context.Context, id string) (string, error)
	Save(ctx context.Context, url URL, userID string) (string, error)
	SaveBatch(ctx context.Context, requestBatch models.PayloadBatch, shortURLBatch models.ResponseBodyBatch, userID string) error
	GetURLs(ctx context.Context, userID string) (models.URLSResponseArr, error)
}

type URL struct {
	OriginalURL string
	ShortURL    string
	GeneratedID string
}
