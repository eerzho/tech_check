package request

type (
	QuestionCreate struct {
		Text       string `json:"text" validate:"required,min=3,max=200"`
		Grade      string `json:"grade" validate:"required,oneof=junior middle senior"`
		CategoryID string `json:"category_id" validate:"required,mongodb"`
	}

	QuestionUpdate struct {
		Text  string `json:"text" validate:"required,min=3,max=200"`
		Grade string `json:"grade" validate:"required,oneof=junior middle senior"`
	}
)
