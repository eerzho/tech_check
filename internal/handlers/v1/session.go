package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handlers/v1/middlewares"
	"tech_check/internal/handlers/v1/requests"
	"tech_check/internal/handlers/v1/responses"
)

type session struct {
	sessionSrvc SessionSrvc
}

func newSession(
	mux *http.ServeMux,
	authMwr *middlewares.Auth,
	sessionSrvc SessionSrvc,
) {
	s := session{
		sessionSrvc: sessionSrvc,
	}

	mux.HandleFunc(
		Url(http.MethodGet, "/sessions"),
		authMwr.MwrFunc(s.list),
	)

	mux.HandleFunc(
		Url(http.MethodPost, "/sessions"),
		authMwr.MwrFunc(s.create),
	)

	mux.HandleFunc(
		Url(http.MethodPost, "/sessions/{id}/summarize"),
		authMwr.MwrFunc(s.summarize),
	)

	mux.HandleFunc(
		Url(http.MethodPost, "/sessions/{id}/cancel"),
		authMwr.MwrFunc(s.cancel),
	)
}

// @Summary get session list for auth user
// @Tags sessions
// @Security BearerAuth
// @Router /v1/sessions [get]
// @Param pagination[page] query int false "page"
// @Param pagination[count] query int false "count"
// @Produce json
// @Success 200 {object} responses.list{data=[]models.Session,pagination=dto.Pagination}
func (s *session) list(w http.ResponseWriter, r *http.Request) {
	const op = "v1.session.list"

	user, err := requests.GetAuthUser(r)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	search := requests.GetQuerySearch(r)
	sessions, pagination, err := s.sessionSrvc.List(
		r.Context(),
		user,
		search.Pagination.Page,
		search.Pagination.Count,
	)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonList(w, r, sessions, pagination)
}

// @Summary start test session
// @Tags sessions
// @Security BearerAuth
// @Router /v1/sessions [post]
// @Param body body requests.SessionCreate true "session create request"
// @Produce json
// @Success 201 {object} responses.success{data=models.Session}
func (s *session) create(w http.ResponseWriter, r *http.Request) {
	const op = "v1.session.create"

	user, err := requests.GetAuthUser(r)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	var req requests.SessionCreate
	err = requests.ParseBody(r, &req)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	session, err := s.sessionSrvc.Create(
		r.Context(),
		user,
		req.CategoryID,
		req.Grade,
	)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusCreated, session)
}

// @Summary finish the session with summary
// @Tags sessions
// @Security BearerAuth
// @Router /v1/sessions/{id}/summarize [post]
// @Param id path string true "session id"
// @Produce json
// @Success 200 {object} responses.success{data=models.Session}
func (s *session) summarize(w http.ResponseWriter, r *http.Request) {
	const op = "v1.session.summarize"

	user, err := requests.GetAuthUser(r)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	id := r.PathValue("id")
	session, err := s.sessionSrvc.Summarize(r.Context(), user, id)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, session)
}

// @Summary finish the session without summary
// @Tags sessions
// @Security BearerAuth
// @Router /v1/sessions/{id}/cancel [post]
// @Param id path string true "session id"
// @Produce json
// @Success 200 {object} responses.success{data=models.Session}
func (s *session) cancel(w http.ResponseWriter, r *http.Request) {
	const op = "v1.session.cancel"

	user, err := requests.GetAuthUser(r)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	id := r.PathValue("id")
	session, err := s.sessionSrvc.Cancel(r.Context(), user, id)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, session)
}
