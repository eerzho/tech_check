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

func (s *SessionQuestion) List(ctx context.Context, session *model.Session) ([]model.SessionQuestion, error) {
	const op = "srvc.SessionQuestion.List"

	questions, err := s.questionRepo.List(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return questions, nil
}

func (s *SessionQuestion) GetByID(ctx context.Context, session *model.Session, id string) (*model.SessionQuestion, error) {
	const op = "srvc.SessionQuestion.GetByID"

	question, err := s.questionRepo.GetByID(ctx, session, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return question, nil
}

func (s *SessionQuestion) Update(ctx context.Context, session *model.Session, id, answer string) (*model.SessionQuestion, error) {
	const op = "srvc.SessionQuestion.Update"

	question, err := s.GetByID(ctx, session, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	question.Answer = answer
	question.Summary = "TODO: ai summary"
	err = s.questionRepo.Update(ctx, question)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return question, nil
}
