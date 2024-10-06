package srvc

import (
	"context"
	"fmt"
	"tech_check/internal/def"
	"tech_check/internal/model"
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

	if session.UserID != user.ID {
		return nil, fmt.Errorf("%s: %w", op, def.ErrAccessDenied)
	}

	return session, nil
}
