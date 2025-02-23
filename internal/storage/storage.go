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

func (m MemoryStorage) Save(ctx context.Context, url store.URL, userID string) (string, error) {
	m[url.GeneratedID] = url.OriginalURL

	return "", nil
}

func (m MemoryStorage) Get(ctx context.Context, id string) (string, bool, error) {
	var originalURL = m[id]

	return originalURL, false, nil
}

func (m MemoryStorage) SaveBatch(ctx context.Context, requestBatch models.PayloadBatch, shortURLBatch models.ResponseBodyBatch, userID string) error {
	var originalURL string

	for _, req := range requestBatch {
		originalURL = req.OriginalURL
	}

	for _, b := range shortURLBatch {
		m[b.CorrelationID] = originalURL
	}

	return nil
}

func (m MemoryStorage) GetURLs(ctx context.Context, userID string) (models.URLSResponseArr, error) {
	return nil, nil
}

func (m MemoryStorage) DeleteURLs(ctx context.Context, userID string, req []string) error {
	return nil
}
