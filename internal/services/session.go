package services

import (
	"context"
	"fmt"
	"tech_check/internal/constants"
	"tech_check/internal/dto"
	"tech_check/internal/models"
	"time"
)

type Session struct {
	count                  int
	sessionRepository      SessionRepository
	categoryService        CategoryService
	questionService        QuestionService
	sessionQuestionService SessionQuestionService
}

func NewSession(
	sessionRepository SessionRepository,
	categoryService CategoryService,
	questionService QuestionService,
	sessionQuestionService SessionQuestionService,
) *Session {
	return &Session{
		count:                  10,
		sessionRepository:      sessionRepository,
		categoryService:        categoryService,
		questionService:        questionService,
		sessionQuestionService: sessionQuestionService,
	}
}

func (s *Session) List(ctx context.Context, user *models.User, page, count int) ([]models.Session, *dto.Pagination, error) {
	const op = "services.Session.List"

	sessions, pagination, err := s.sessionRepository.List(ctx, user, page, count)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return sessions, pagination, nil
}

func (s *Session) Create(ctx context.Context, user *models.User, categoryID, grade string) (*models.Session, error) {
	const op = "services.Session.Create"

	exists, err := s.sessionRepository.ExistsActive(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if exists {
		return nil, fmt.Errorf("%s: %w", op, constants.ErrUserHasActiveSession)
	}

	category, err := s.categoryService.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	questions, err := s.questionService.GetRandom(ctx, category, grade, s.count)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(questions) != s.count {
		return nil, fmt.Errorf("%s: %w", op, constants.ErrQuestionNotEnough)
	}

	gradeObj, err := constants.ValidateGradeName(grade)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	session := models.Session{
		UserID:     user.ID,
		CategoryID: category.ID,
		Grade:      gradeObj,
	}
	err = s.sessionRepository.Create(ctx, &session)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, question := range questions {
		_, err := s.sessionQuestionService.Create(ctx, &session, question.Text)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &session, nil
}

func (s *Session) GetByID(ctx context.Context, user *models.User, id string) (*models.Session, error) {
	const op = "services.Session.GetByID"

	session, err := s.sessionRepository.GetByID(ctx, user, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (s *Session) Summarize(ctx context.Context, user *models.User, id string) (*models.Session, error) {
	const op = "services.Session.Summarize"

	session, err := s.GetByID(ctx, user, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()
	session.Summary = "TODO: ai summary"
	session.FinishedAt = &now
	err = s.sessionRepository.Update(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (s *Session) Cancel(ctx context.Context, user *models.User, id string) (*models.Session, error) {
	const op = "services.Session.Cancel"

	session, err := s.GetByID(ctx, user, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()
	session.FinishedAt = &now
	err = s.sessionRepository.Update(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}
