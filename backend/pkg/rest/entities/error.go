package entities

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
