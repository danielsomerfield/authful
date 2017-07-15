package wire

type ErrorsResponse struct {
	Error            []Error `json:"error"`
}

type Error struct {
	Status	string `json:"status"`
	Detail	string `json:"detail"`
}
