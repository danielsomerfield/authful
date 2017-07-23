package wire

type ErrorsResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Status    int       `json:"status"`
	Detail    string    `json:"detail"`
	ErrorType string    `json:"errorType"`
	ErrorURI  string    `json:"errorURI"`
}
