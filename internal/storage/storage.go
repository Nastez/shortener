package storage

import (
	"context"

	"github.com/Nastez/shortener/internal/store"
)

type MemoryStorage map[string]string

func (m MemoryStorage) Bootstrap(ctx context.Context) error {
	return nil
}

func (m MemoryStorage) Save(ctx context.Context, url store.URL) error {
	m[url.GeneratedID] = url.OriginalURL

	return nil
}

func (m MemoryStorage) Get(ctx context.Context, id string) (string, error) {
	var originalURL = m[id]

	return originalURL, nil
}
