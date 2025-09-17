package model

type (
	ErrorResponse struct {
		Error   string  `json:"error"`
		Message *string `json:"message"`
	}
)

func NewErrorResponse(error string, message ...string) *ErrorResponse {
	e := &ErrorResponse{Error: error}

	if len(message) > 0 {
		e.Message = &message[0]
	}

	return e
}
