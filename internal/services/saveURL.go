package services

//
//import (
//	"context"
//	"fmt"
//	"github.com/Nastez/shortener/internal/store"
//	"github.com/Nastez/shortener/utils"
//)
//
//type SaveService struct {
//	store    store.Store
//	baseAddr string
//}
//
//func NewSaveService(s store.Store, baseAddr string) *SaveService {
//	return &SaveService{store: s, baseAddr: baseAddr}
//}
//
//func (s *SaveService) SaveURL(ctx context.Context, originalURL string) (string, string, error) {
//	generatedID := utils.GenerateID()
//	var shortURL string
//
//	shortURL = s.baseAddr + "/" + generatedID
//
//	fmt.Println(originalURL)
//
//	oldShortURL, err := s.store.Save(ctx, store.URL{
//		OriginalURL: originalURL,
//		ShortURL:    shortURL,
//		GeneratedID: generatedID,
//	})
//
//	fmt.Println(oldShortURL, shortURL)
//
//	return oldShortURL, shortURL, err
//}
