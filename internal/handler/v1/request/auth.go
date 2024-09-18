package request

type (
	Login struct {
		Email    string `json:"email" validate:"required,email,max=50"`
		Password string `json:"password" validate:"required,min=8,max=50"`
	}

	Refresh struct {
		AToken string `json:"access_token" validate:"required"`
		RToken string `json:"refresh_token" validate:"required"`
	}
)
