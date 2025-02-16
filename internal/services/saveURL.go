package services

import (
	"context"

	"github.com/Nastez/shortener/internal/store"
	"github.com/Nastez/shortener/utils"
)

func SaveURL(ctx context.Context, baseAddr string, storage store.Store, originalURL string, userID string) (string, string, error) {
	generatedID := utils.GenerateID()
	shortURL := baseAddr + "/" + generatedID

	oldShortURL, err := storage.Save(ctx, store.URL{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		GeneratedID: generatedID,
	}, userID)

	return oldShortURL, shortURL, err
}
