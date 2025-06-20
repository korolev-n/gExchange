package httptransport

type RatesResponse struct {
	Rates map[string]float64 `json:"rates"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
