package helper

// BaseErrorResponse defines the general error response structure
type BaseErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}
