package services

import (
	"context"
	"fmt"
	"tech_check/internal/constants"
	"tech_check/internal/dto"
	"tech_check/internal/models"
)

type Question struct {
	questionRepo QuestionRepo
	categorySrvc CategorySrvc
}

func NewQuestion(
	questionRepo QuestionRepo,
	categorySrvc CategorySrvc,
) *Question {
	return &Question{
		questionRepo: questionRepo,
		categorySrvc: categorySrvc,
	}
}

func (q *Question) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Question, *dto.Pagination, error) {
	const op = "srvq.Question.List"

	questions, pagination, err := q.questionRepo.List(ctx, page, count, filters, sorts)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return questions, pagination, nil
}

func (q *Question) Create(ctx context.Context, text, grade, categoryID string) (*models.Question, error) {
	const op = "srvq.Question.Create"

	gradeObj, err := constants.ValidateGradeName(grade)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	category, err := q.categorySrvc.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	question := models.Question{
		Text:       text,
		Grade:      gradeObj,
		CategoryID: category.ID,
	}
	err = q.questionRepo.Create(ctx, &question)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &question, nil
}

func (q *Question) GetByID(ctx context.Context, id string) (*models.Question, error) {
	const op = "srvq.Question.GetByID"

	question, err := q.questionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return question, nil
}

func (q *Question) Update(ctx context.Context, id, text, grade string) (*models.Question, error) {
	const op = "srvq.Question.Update"

	gradeObj, err := constants.ValidateGradeName(grade)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	question, err := q.questionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	question.Text = text
	question.Grade = gradeObj
	err = q.questionRepo.Update(ctx, question)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return question, nil
}

func (q *Question) Delete(ctx context.Context, id string) error {
	const op = "services.Question.Delete"

	err := q.questionRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (q *Question) GetRandom(ctx context.Context, category *models.Category, grade string, count int) ([]models.Question, error) {
	const op = "services.Question.GetRandom"

	gradeObj, err := constants.ValidateGradeName(grade)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	questions, err := q.questionRepo.GetRandom(ctx, category, gradeObj, count)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return questions, nil
}
