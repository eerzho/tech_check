package srvc

import (
	"context"
	"fmt"
	"tech_check/internal/model"
)

type SessionQuestion struct {
	questionRepo SessionQuestionRepo
}

func NewSessionQuestion(
	questionRepo SessionQuestionRepo,
) *SessionQuestion {
	return &SessionQuestion{
		questionRepo: questionRepo,
	}
}

func (s *SessionQuestion) Create(ctx context.Context, session *model.Session, text string) (*model.SessionQuestion, error) {
	const op = "srvc.SessionQuestion.Create"

	question := model.SessionQuestion{
		SessionID: session.ID,
		Text:      text,
	}
	err := s.questionRepo.Create(ctx, &question)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &question, nil
}
