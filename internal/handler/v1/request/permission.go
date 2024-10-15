package request

type (
	PermissionCreate struct {
		Name string `json:"name" validate:"required,min=5,max=50"`
	}

	PermissionUpdate struct {
		Name string `json:"name" validate:"required,min=5,max=50"`
	}
)
