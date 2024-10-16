package services

import (
	"context"
	"fmt"
	"tech_check/internal/constants"
	"tech_check/internal/dto"
	"tech_check/internal/models"
)

type Question struct {
	questionRepository QuestionRepository
	categoryService    CategoryService
}

func NewQuestion(
	questionRepository QuestionRepository,
	categoryService CategoryService,
) *Question {
	return &Question{
		questionRepository: questionRepository,
		categoryService:    categoryService,
	}
}

func (q *Question) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Question, *dto.Pagination, error) {
	const op = "srvq.Question.List"

	questions, pagination, err := q.questionRepository.List(ctx, page, count, filters, sorts)
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

	category, err := q.categoryService.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	question := models.Question{
		Text:       text,
		Grade:      gradeObj,
		CategoryID: category.ID,
	}
	err = q.questionRepository.Create(ctx, &question)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &question, nil
}

func (q *Question) GetByID(ctx context.Context, id string) (*models.Question, error) {
	const op = "srvq.Question.GetByID"

	question, err := q.questionRepository.GetByID(ctx, id)
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

	question, err := q.questionRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	question.Text = text
	question.Grade = gradeObj
	err = q.questionRepository.Update(ctx, question)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return question, nil
}

func (q *Question) Delete(ctx context.Context, id string) error {
	const op = "services.Question.Delete"

	err := q.questionRepository.Delete(ctx, id)
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

	questions, err := q.questionRepository.GetRandom(ctx, category, gradeObj, count)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return questions, nil
}
