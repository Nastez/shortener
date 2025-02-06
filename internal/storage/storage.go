//go:generate mockgen -source=storage.go -destination=mocks/mocks.go
package storage

//type MemoryStorage map[string]string
//
//type URLStorage interface {
//	Save(originalURL string, generatedID string)
//	Get(urlID string) string
//}
//
//func (m MemoryStorage) Save(originalURL string, generatedID string) {
//	m[generatedID] = originalURL
//}
//
//func (m MemoryStorage) Get(urlID string) string {
//	var originalURL = m[urlID]
//
//	return originalURL
//}
