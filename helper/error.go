package helper

// BaseErrorResponse defines the general error response structure
type BaseErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func (b *BaseErrorResponse) Error() string {
	if b == nil {
		return "unknown error"
	}
	return b.Message
}
