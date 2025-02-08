package models

type Request struct {
	URL string `json:"url"`
}

type PayloadBatch []RequestBatch

type RequestBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
