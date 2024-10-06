package request

type (
	SessionCreate struct {
		CategoryID string `json:"category_id" validate:"required,mongodb"`
		Grade      string `json:"grade" validate:"required,oneof=junior middle senior"`
	}
)
