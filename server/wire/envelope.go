package wire

type ResponseEnvelope struct {
	Data interface{}        `json:"data,omitempty"`
}

type RequestEnvelope struct {
	Command interface{}        `json:"command,omitempty"`
}