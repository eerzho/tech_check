package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handlers/v1/middlewares"
	"tech_check/internal/handlers/v1/requests"
	"tech_check/internal/handlers/v1/responses"
)

type question struct {
	questionService QuestionService
}

func newQuestion(
	mux *http.ServeMux,
	authMwr *middlewares.Auth,
	permissionMwr *middlewares.Permission,
	questionService QuestionService,
) {
	q := question{
		questionService: questionService,
	}

	mux.HandleFunc(
		Url(http.MethodGet, "/questions"),
		authMwr.MwrFunc(q.list),
	)

	mux.HandleFunc(
		Url(http.MethodPost, "/questions"),
		authMwr.MwrFunc(permissionMwr.MwrFunc(q.create, "question-create")),
	)

	mux.HandleFunc(
		Url(http.MethodGet, "/questions/{id}"),
		authMwr.MwrFunc(q.show),
	)

	mux.HandleFunc(
		Url(http.MethodPatch, "/questions/{id}"),
		authMwr.MwrFunc(permissionMwr.MwrFunc(q.update, "question-edit")),
	)

	mux.HandleFunc(
		Url(http.MethodDelete, "/questions/{id}"),
		authMwr.MwrFunc(permissionMwr.MwrFunc(q.delete, "question-delete")),
	)
}

// @Summary questions list
// @Tags questions
// @Security BearerAuth
// @Router /v1/questions [get]
// @Param pagination[page] query int false "page"
// @Param pagination[count] query int false "count"
// @Param sorts[created_at] query string false "created_at" Enums(asc, desc)
// @Param sorts[updated_at] query string false "updated_at" Enums(asc, desc)
// @Param filters[text] query string false "text"
// @Param filters[grade] query string false "grade" Enums(junior, middle, senior)
// @Produce json
// @Success 200 {object} responses.list{data=[]models.Question,pagination=dto.Pagination}
func (q *question) list(w http.ResponseWriter, r *http.Request) {
	const op = "v1.question.list"

	search := requests.GetQuerySearch(r)
	questions, pagination, err := q.questionService.List(
		r.Context(),
		search.Pagination.Page,
		search.Pagination.Count,
		search.Filters,
		search.Sorts,
	)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonList(w, r, questions, pagination)
}

// @Summary create question
// @Tags questions
// @Security BearerAuth
// @Router /v1/questions [post]
// @Accept json
// @Param body body requests.QuestionCreate true "question create request"
// @Produce json
// @Success 201 {object} responses.success{data=models.Question}
func (q *question) create(w http.ResponseWriter, r *http.Request) {
	const op = "v1.question.create"

	var req requests.QuestionCreate
	err := requests.ParseBody(r, &req)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	question, err := q.questionService.Create(
		r.Context(),
		req.Text,
		req.Grade,
		req.CategoryID,
	)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusCreated, question)
}

// @Summary get question by id
// @Tags questions
// @Security BearerAuth
// @Router /v1/questions/{id} [get]
// @Param id path string true "question id"
// @Produce json
// @Success 200 {object} responses.success{data=models.Question}
func (q *question) show(w http.ResponseWriter, r *http.Request) {
	const op = "v1.question.show"

	id := r.PathValue("id")
	question, err := q.questionService.GetByID(r.Context(), id)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, question)
}

// @Summary update profile
// @Tags questions
// @Security BearerAuth
// @Router /v1/questions/{id} [patch]
// @Accept json
// @Param id path string true "question id"
// @Param body body requests.QuestionUpdate true "question update request"
// @Produce json
// @Success 200 {object} responses.success{data=models.Question}
func (q *question) update(w http.ResponseWriter, r *http.Request) {
	const op = "v1.question.update"

	id := r.PathValue("id")
	var req requests.QuestionUpdate

	err := requests.ParseBody(r, &req)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	question, err := q.questionService.Update(r.Context(), id, req.Text, req.Grade)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, question)
}

// @Summary delete question by id
// @Tags questions
// @Security BearerAuth
// @Router /v1/questions/{id} [delete]
// @Param id path string true "question id"
// @Success 204
func (q *question) delete(w http.ResponseWriter, r *http.Request) {
	const op = "v1.question.delete"

	id := r.PathValue("id")

	err := q.questionService.Delete(r.Context(), id)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusNoContent, nil)
}
