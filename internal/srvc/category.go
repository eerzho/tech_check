package srvc

import (
	"context"
	"fmt"
	"tech_check/internal/dto"
	"tech_check/internal/model"

	"github.com/gosimple/slug"
)

type Category struct {
	categoryRepo CategoryRepo
}

func NewCategory(categoryRepo CategoryRepo) *Category {
	return &Category{
		categoryRepo: categoryRepo,
	}
}

func (c *Category) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Category, *dto.Pagination, error) {
	const op = "srvc.Category.List"

	categories, pagination, err := c.categoryRepo.List(ctx, page, count, filters, sorts)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return categories, pagination, nil
}

func (c *Category) Create(ctx context.Context, name, description string) (*model.Category, error) {
	const op = "srvc.Category.Create"

	slug := slug.Make(name)
	count, err := c.categoryRepo.CountBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if count > 0 {
		slug = fmt.Sprintf("%s-%d", slug, count+1)
	}

	category := model.Category{
		Name:        name,
		Slug:        slug,
		Description: description,
	}
	err = c.categoryRepo.Create(ctx, &category)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &category, nil
}

func (c *Category) GetByID(ctx context.Context, id string) (*model.Category, error) {
	const op = "srvc.Category.GetByID"

	category, err := c.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return category, nil
}

func (c *Category) Update(ctx context.Context, id, name, description string) (*model.Category, error) {
	const op = "srvc.Category.Update"

	category, err := c.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	category.Name = name
	category.Description = description
	err = c.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return category, nil
}

func (c *Category) Delete(ctx context.Context, id string) error {
	const op = "srvc.Category.Delete"

	err := c.categoryRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
