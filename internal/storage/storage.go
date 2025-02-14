package storage

import (
	"context"

	"github.com/Nastez/shortener/internal/app/models"
	"github.com/Nastez/shortener/internal/store"
)

type MemoryStorage map[string]string

func New() *MemoryStorage {
	return &MemoryStorage{}
}

func (m MemoryStorage) Bootstrap(ctx context.Context) error {
	return nil
}

func (m MemoryStorage) Save(ctx context.Context, url store.URL) (string, error) {
	m[url.GeneratedID] = url.OriginalURL

	return "", nil
}

func (m MemoryStorage) Get(ctx context.Context, id string) (string, error) {
	var originalURL = m[id]

	return originalURL, nil
}

func (m MemoryStorage) SaveBatch(ctx context.Context, requestBatch models.PayloadBatch, shortURLBatch models.ResponseBodyBatch) error {
	var originalURL string

	for _, req := range requestBatch {
		originalURL = req.OriginalURL
	}

	for _, b := range shortURLBatch {
		m[b.CorrelationID] = originalURL
	}

	return nil
}
