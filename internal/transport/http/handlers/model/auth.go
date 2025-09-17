package model

type (
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	TokenResponse struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}
)
