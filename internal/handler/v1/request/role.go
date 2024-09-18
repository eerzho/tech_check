package request

type RoleCreate struct {
	Name string `json:"name" validate:"required,min=5,max=50"`
}

type RoleUpdate struct {
	Name string `json:"name" validate:"required,min=5,max=50"`
}
