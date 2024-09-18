package request

type (
	CategoryCreate struct {
		Name        string `json:"name" validate:"required,min=1,max=50"`
		Description string `json:"description" validate:"required,min=5,max=500"`
	}

	CategoryUpdate struct {
		Name        string `json:"name" validate:"required,min=1,max=50"`
		Description string `json:"description" validate:"required,min=5,max=500"`
	}
)
