package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

type category struct {
	rp           *request.Parser
	rb           *response.Builder
	categorySrvc CategorySrvc
}

func newCategory(
	mux *http.ServeMux,
	prefix string,
	rp *request.Parser,
	rb *response.Builder,
	authMwr *mwr.Auth,
	categorySrvc CategorySrvc,
) {
	prefix += "/categories"
	c := category{
		rp:           rp,
		rb:           rb,
		categorySrvc: categorySrvc,
	}

	mux.HandleFunc("GET "+prefix, authMwr.MwrFunc(c.list))
	mux.HandleFunc("POST "+prefix, authMwr.MwrFunc(c.create))
	mux.HandleFunc("GET "+prefix+"/{id}", authMwr.MwrFunc(c.show))
	mux.HandleFunc("PATCH "+prefix+"/{id}", authMwr.MwrFunc(c.update))
	mux.HandleFunc("DELETE "+prefix+"/{id}", authMwr.MwrFunc(c.delete))
}

// @Summary categories list
// @Tags categories
// @Security BearerAuth
// @Router /v1/categories [get]
// @Param pagination[page] query int false "page"
// @Param pagination[count] query int false "count"
// @Param sorts[created_at] query string false "created_at" Enums(asc, desc)
// @Param sorts[updated_at] query string false "updated_at" Enums(asc, desc)
// @Param sorts[name] query string false "name" Enums(asc, desc)
// @Param sorts[slug] query string false "slug" Enums(asc, desc)
// @Param filters[name] query string false "name"
// @Param filters[slug] query string false "slug"
// @Param filters[description] query string false "description"
// @Produce json
// @Success 200 {object} response.list{data=[]model.Category,pagination=dto.Pagination}
func (c *category) list(w http.ResponseWriter, r *http.Request) {
	const op = "v1.category.list"

	search := c.rp.GetQuerySearch(r)
	categories, pagination, err := c.categorySrvc.List(
		r.Context(),
		search.Pagination.Page,
		search.Pagination.Count,
		search.Filters,
		search.Sorts,
	)
	if err != nil {
		c.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	c.rb.JsonList(w, r, categories, pagination)
}

// @Summary create category
// @Tags categories
// @Security BearerAuth
// @Router /v1/categories [post]
// @Accept json
// @Param body body request.CategoryCreate true "category create request"
// @Produce json
// @Success 201 {object} response.success{data=model.Category}
func (c *category) create(w http.ResponseWriter, r *http.Request) {
	const op = "v1.category.create"

	var req request.CategoryCreate
	err := c.rp.ParseBody(r, &req)
	if err != nil {
		c.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	category, err := c.categorySrvc.Create(
		r.Context(),
		req.Name,
		req.Description,
	)
	if err != nil {
		c.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	c.rb.JsonSuccess(w, r, http.StatusCreated, category)
}

// @Summary get category by id
// @Tags categories
// @Security BearerAuth
// @Router /v1/categories/{id} [get]
// @Param id path string true "category id"
// @Produce json
// @Success 200 {object} response.success{data=model.Category}
func (c *category) show(w http.ResponseWriter, r *http.Request) {
	const op = "v1.category.show"

	id := r.PathValue("id")
	category, err := c.categorySrvc.GetByID(r.Context(), id)
	if err != nil {
		c.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	c.rb.JsonSuccess(w, r, http.StatusOK, category)
}

// @Summary update profile
// @Tags categories
// @Security BearerAuth
// @Router /v1/categories/{id} [patch]
// @Accept json
// @Param id path string true "category id"
// @Param body body request.CategoryUpdate true "category update request"
// @Produce json
// @Success 200 {object} response.success{data=model.Category}
func (c *category) update(w http.ResponseWriter, r *http.Request) {
	const op = "v1.category.update"

	id := r.PathValue("id")
	var req request.CategoryUpdate

	err := c.rp.ParseBody(r, &req)
	if err != nil {
		c.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	category, err := c.categorySrvc.Update(r.Context(), id, req.Name, req.Description)
	if err != nil {
		c.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	c.rb.JsonSuccess(w, r, http.StatusOK, category)
}

// @Summary delete category by id
// @Tags categories
// @Security BearerAuth
// @Router /v1/categories/{id} [delete]
// @Param id path string true "category id"
// @Success 204
func (c *category) delete(w http.ResponseWriter, r *http.Request) {
	const op = "v1.category.delete"

	id := r.PathValue("id")

	err := c.categorySrvc.Delete(r.Context(), id)
	if err != nil {
		c.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	c.rb.JsonSuccess(w, r, http.StatusNoContent, nil)
}
