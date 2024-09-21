package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

type question struct {
	questionSrvc QuestionSrvc
}

func newQuestion(
	mux *http.ServeMux,
	prefix string,
	authMwr *mwr.Auth,
	questionSrvc QuestionSrvc,
) {
	prefix += "/questions"
	q := question{
		questionSrvc: questionSrvc,
	}

	mux.HandleFunc("GET "+prefix, authMwr.MwrFunc(q.list))
	mux.HandleFunc("POST "+prefix, authMwr.MwrFunc(q.create))
	mux.HandleFunc("GET "+prefix+"/{id}", authMwr.MwrFunc(q.show))
	mux.HandleFunc("PATCH "+prefix+"/{id}", authMwr.MwrFunc(q.update))
	mux.HandleFunc("DELETE "+prefix+"/{id}", authMwr.MwrFunc(q.delete))
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
// @Produce json
// @Success 200 {object} response.list{data=[]model.Question,pagination=dto.Pagination}
func (q *question) list(w http.ResponseWriter, r *http.Request) {
	const op = "v1.question.list"

	search := request.GetQuerySearch(r)
	questions, pagination, err := q.questionSrvc.List(
		r.Context(),
		search.Pagination.Page,
		search.Pagination.Count,
		search.Filters,
		search.Sorts,
	)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonList(w, r, questions, pagination)
}

// @Summary create question
// @Tags questions
// @Security BearerAuth
// @Router /v1/questions [post]
// @Accept json
// @Param body body request.QuestionCreate true "question create request"
// @Produce json
// @Success 201 {object} response.success{data=model.Question}
func (q *question) create(w http.ResponseWriter, r *http.Request) {
	const op = "v1.question.create"

	var req request.QuestionCreate
	err := request.ParseBody(r, &req)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	question, err := q.questionSrvc.Create(
		r.Context(),
		req.Text,
		req.CategoryID,
	)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusCreated, question)
}

// @Summary get question by id
// @Tags questions
// @Security BearerAuth
// @Router /v1/questions/{id} [get]
// @Param id path string true "question id"
// @Produce json
// @Success 200 {object} response.success{data=model.Question}
func (q *question) show(w http.ResponseWriter, r *http.Request) {
	const op = "v1.question.show"

	id := r.PathValue("id")
	question, err := q.questionSrvc.GetByID(r.Context(), id)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusOK, question)
}

// @Summary update profile
// @Tags questions
// @Security BearerAuth
// @Router /v1/questions/{id} [patch]
// @Accept json
// @Param id path string true "question id"
// @Param body body request.QuestionUpdate true "question update request"
// @Produce json
// @Success 200 {object} response.success{data=model.Question}
func (q *question) update(w http.ResponseWriter, r *http.Request) {
	const op = "v1.question.update"

	id := r.PathValue("id")
	var req request.QuestionUpdate

	err := request.ParseBody(r, &req)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	question, err := q.questionSrvc.Update(r.Context(), id, req.Text)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusOK, question)
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

	err := q.questionSrvc.Delete(r.Context(), id)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusNoContent, nil)
}
