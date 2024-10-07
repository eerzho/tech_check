package srvc

import (
	"context"
	"fmt"
	"tech_check/internal/def"
	"tech_check/internal/dto"
	"tech_check/internal/model"
	"time"
)

type Session struct {
	count               int
	sessionRepo         SessionRepo
	categorySrvc        CategorySrvc
	questionSrvc        QuestionSrvc
	sessionQuestionSrvc SessionQuestionSrvc
}

func NewSession(
	sessionRepo SessionRepo,
	categorySrvc CategorySrvc,
	questionSrvc QuestionSrvc,
	sessionQuestionSrvc SessionQuestionSrvc,
) *Session {
	return &Session{
		count:               10,
		sessionRepo:         sessionRepo,
		categorySrvc:        categorySrvc,
		questionSrvc:        questionSrvc,
		sessionQuestionSrvc: sessionQuestionSrvc,
	}
}

func (s *Session) List(ctx context.Context, user *model.User, page, count int) ([]model.Session, *dto.Pagination, error) {
	const op = "srvc.Session.List"

	sessions, pagination, err := s.sessionRepo.List(ctx, user, page, count)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return sessions, pagination, nil
}

func (s *Session) Create(ctx context.Context, user *model.User, categoryID, grade string) (*model.Session, error) {
	const op = "srvc.Session.Create"

	exists, err := s.sessionRepo.IsExistsActive(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if exists {
		return nil, fmt.Errorf("%s: %w", op, def.ErrUserHasActiveSession)
	}

	category, err := s.categorySrvc.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	questions, err := s.questionSrvc.GetRandom(ctx, category, grade, s.count)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(questions) != s.count {
		return nil, fmt.Errorf("%s: %w", op, def.ErrQuestionNotEnough)
	}

	gradeObj, err := def.ValidateGradeName(grade)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	session := model.Session{
		UserID:     user.ID,
		CategoryID: category.ID,
		Grade:      gradeObj,
	}
	err = s.sessionRepo.Create(ctx, &session)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, question := range questions {
		_, err := s.sessionQuestionSrvc.Create(ctx, &session, question.Text)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &session, nil
}

func (s *Session) GetByID(ctx context.Context, user *model.User, id string) (*model.Session, error) {
	const op = "srvc.Session.GetByID"

	session, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if session.FinishedAt != nil {
		return nil, fmt.Errorf("%s: %w", op, def.ErrSessionFinished)
	}

	if session.UserID != user.ID {
		return nil, fmt.Errorf("%s: %w", op, def.ErrAccessDenied)
	}

	return session, nil
}

func (s *Session) Summarize(ctx context.Context, user *model.User, id string) (*model.Session, error) {
	const op = "srvc.Session.Summarize"

	session, err := s.GetByID(ctx, user, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()
	session.Summary = "TODO: ai summary"
	session.FinishedAt = &now
	err = s.sessionRepo.Update(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (s *Session) Cancel(ctx context.Context, user *model.User, id string) (*model.Session, error) {
	const op = "srvc.Session.Cancel"

	session,  err := s.GetByID(ctx, user, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()
	session.FinishedAt = &now
	err = s.sessionRepo.Update(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}
