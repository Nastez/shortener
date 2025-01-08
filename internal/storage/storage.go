package storage

type MemoryStorage map[string]string

type URLStorage interface {
	SaveOriginalURL(originalURL string, generatedID string)
	GetOriginalURL(urlID string) string
}

func (m MemoryStorage) SaveOriginalURL(originalURL string, generatedID string) {
	m[generatedID] = originalURL
}

func (m MemoryStorage) GetOriginalURL(urlID string) string {
	var originalURL = m[urlID]

	return originalURL
}
