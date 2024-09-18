package request

type (
	UserCreate struct {
		Email    string `json:"email" validate:"required,email,max=50"`
		Name     string `json:"name" validate:"required,max=50"`
		Password string `json:"password" validate:"required,min=8,max=50"`
	}

	UserUpdate struct {
		Name string `json:"name" validate:"required,max=50"`
	}
)
