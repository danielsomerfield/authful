package wire

type ErrorsResponse struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Status    int       `json:"status"`
	Detail    string    `json:"detail"`
	ErrorType string    `json:"errorType"`
	ErrorURI  string    `json:"errorURI"`
}
