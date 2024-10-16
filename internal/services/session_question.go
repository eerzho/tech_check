package services

import (
	"context"
	"fmt"
	"tech_check/internal/models"
)

type SessionQuestion struct {
	questionRepository SessionQuestionRepository
}

func NewSessionQuestion(
	questionRepository SessionQuestionRepository,
) *SessionQuestion {
	return &SessionQuestion{
		questionRepository: questionRepository,
	}
}

func (s *SessionQuestion) Create(ctx context.Context, session *models.Session, text string) (*models.SessionQuestion, error) {
	const op = "services.SessionQuestion.Create"

	question := models.SessionQuestion{
		SessionID: session.ID,
		Text:      text,
	}
	err := s.questionRepository.Create(ctx, &question)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &question, nil
}

func (s *SessionQuestion) List(ctx context.Context, session *models.Session) ([]models.SessionQuestion, error) {
	const op = "services.SessionQuestion.List"

	questions, err := s.questionRepository.List(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return questions, nil
}

func (s *SessionQuestion) GetByID(ctx context.Context, session *models.Session, id string) (*models.SessionQuestion, error) {
	const op = "services.SessionQuestion.GetByID"

	question, err := s.questionRepository.GetByID(ctx, session, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return question, nil
}

func (s *SessionQuestion) Update(ctx context.Context, session *models.Session, id, answer string) (*models.SessionQuestion, error) {
	const op = "services.SessionQuestion.Update"

	question, err := s.GetByID(ctx, session, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	question.Answer = answer
	question.Summary = "TODO: ai summary"
	err = s.questionRepository.Update(ctx, question)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return question, nil
}
