package models

type Response struct {
	Result string `json:"result"`
}

type ResponseBodyBatch []ResponseBatch

type ResponseBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
