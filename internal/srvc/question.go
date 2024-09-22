package srvc

import (
	"context"
	"fmt"
	"tech_check/internal/def"
	"tech_check/internal/dto"
	"tech_check/internal/model"
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

func (q *Question) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Question, *dto.Pagination, error) {
	const op = "srvq.Question.List"

	questions, pagination, err := q.questionRepo.List(ctx, page, count, filters, sorts)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return questions, pagination, nil
}

func (q *Question) Create(ctx context.Context, text, grade, categoryID string) (*model.Question, error) {
	const op = "srvq.Question.Create"

	gradeObj, err := def.ValidateGradeName(grade)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	category, err := q.categorySrvc.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	question := model.Question{
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

func (q *Question) GetByID(ctx context.Context, id string) (*model.Question, error) {
	const op = "srvq.Question.GetByID"

	question, err := q.questionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return question, nil
}

func (q *Question) Update(ctx context.Context, id, text, grade string) (*model.Question, error) {
	const op = "srvq.Question.Update"

	gradeObj, err := def.ValidateGradeName(grade)
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
	const op = "srvq.Question.Delete"

	err := q.questionRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
