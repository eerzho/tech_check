package srvc

import (
	"context"
	"fmt"
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

func (q *Question) Create(ctx context.Context, text, categoryID string) (*model.Question, error) {
	const op = "srvq.Question.Create"

	category, err := q.categorySrvc.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	question := model.Question{
		Text:       text,
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

func (q *Question) Update(ctx context.Context, id, text string) (*model.Question, error) {
	const op = "srvq.Question.Update"

	question, err := q.questionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	question.Text = text
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
