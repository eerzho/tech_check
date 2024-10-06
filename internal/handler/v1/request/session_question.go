package request

type SessionQuestionUpdate struct {
	Answer string `json:"answer" validate:"required,min=1,max=500"`
}
