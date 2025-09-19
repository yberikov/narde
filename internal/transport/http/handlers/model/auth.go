package model

type (
	RegisterRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	TokenResponse struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}
)
