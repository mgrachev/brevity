package api

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessfulReponse struct {
	Message string `json:"message"`
	Link    string `json:"link"`
}

type FailureResponse struct {
	Field string `json:"field"`
	Error string `json:"error"`
}
