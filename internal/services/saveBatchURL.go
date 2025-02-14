package services

import (
	"context"

	"github.com/Nastez/shortener/internal/app/models"
	"github.com/Nastez/shortener/internal/logger"
	"github.com/Nastez/shortener/internal/store"
)

func SaveBatchURL(ctx context.Context, requestBatch models.PayloadBatch, baseAddr string, storage store.Store) models.ResponseBodyBatch {
	var responseBatch models.ResponseBodyBatch

	for _, request := range requestBatch {
		var response = models.ResponseBatch{
			CorrelationID: request.CorrelationID,
			ShortURL:      baseAddr + "/" + request.CorrelationID,
		}
		responseBatch = append(responseBatch, response)
	}

	if len(responseBatch) > 0 {
		err := storage.SaveBatch(ctx, requestBatch, responseBatch)
		if err != nil {
			logger.Log.Info("can't save batch in store")
			return nil
		}
	} else {
		logger.Log.Info("responseBatch is empty")
	}

	return responseBatch
}
