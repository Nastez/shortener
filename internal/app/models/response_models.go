package models

type Response struct {
	Result string `json:"result"`
}

type ResponseBodyBatch []ResponseBatch

type ResponseBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type URLSResponseArr []URLSResponse

type URLSResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
