package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

type sessionQuestion struct {
	sessionSrvc         SessionSrvc
	sessionQuestionSrvc SessionQuestionSrvc
}

func newSessionQuestion(
	mux *http.ServeMux,
	authMwr *mwr.Auth,
	sessionSrvc SessionSrvc,
	sessionQuestionSrvc SessionQuestionSrvc,
) {
	s := sessionQuestion{
		sessionSrvc:         sessionSrvc,
		sessionQuestionSrvc: sessionQuestionSrvc,
	}

	mux.HandleFunc(
		Url(http.MethodGet, "/sessions/{sessionID}/questions"),
		authMwr.MwrFunc(s.list),
	)

	mux.HandleFunc(
		Url(http.MethodGet, "/sessions/{sessionID}/questions/{id}"),
		authMwr.MwrFunc(s.show),
	)

	mux.HandleFunc(
		Url(http.MethodPatch, "/sessions/{sessionID}/questions/{id}"),
		authMwr.MwrFunc(s.update),
	)
}

// @Summary get session questions
// @Tags sessionQuestions
// @Security BearerAuth
// @Router /v1/sessions/{sessionID}/questions [get]
// @Param sessionID path string true "session id"
// @Produce json
// @Success 200 {object} response.success{data=[]model.SessionQuestion}
func (s *sessionQuestion) list(w http.ResponseWriter, r *http.Request) {
	const op = "v1.sessionQuestion.list"

	user, err := request.GetAuthUser(r)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	sessionID := r.PathValue("sessionID")
	session, err := s.sessionSrvc.GetByID(r.Context(), user, sessionID)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	questions, err := s.sessionQuestionSrvc.List(r.Context(), session)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
	}

	response.JsonSuccess(w, r, http.StatusOK, questions)
}

// @Summary get session question by id
// @Tags sessionQuestions
// @Security BearerAuth
// @Router /v1/sessions/{sessionID}/questions/{id} [get]
// @Param sessionID path string true "session id"
// @Param id path string true "session question id"
// @Produce json
// @Success 200 {object} response.success{data=model.SessionQuestion}
func (s *sessionQuestion) show(w http.ResponseWriter, r *http.Request) {
	const op = "v1.sessionQuestion.show"

	user, err := request.GetAuthUser(r)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	sessionID := r.PathValue("sessionID")
	session, err := s.sessionSrvc.GetByID(r.Context(), user, sessionID)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	id := r.PathValue("id")
	question, err := s.sessionQuestionSrvc.GetByID(r.Context(), session, id)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusOK, question)
}

// @Summary answer the question
// @Tags sessionQuestions
// @Security BearerAuth
// @Router /v1/sessions/{sessionID}/questions/{id} [patch]
// @Accept json
// @Param sessionID path string true "session id"
// @Param id path string true "session question id"
// @Param body body request.SessionQuestionUpdate true "answer the question request"
// @Produce json
// @Success 200 {object} response.success{data=model.SessionQuestion}
func (s *sessionQuestion) update(w http.ResponseWriter, r *http.Request) {
	const op = "v1.sessionQuestion.update"

	var req request.SessionQuestionUpdate
	err := request.ParseBody(r, &req)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	user, err := request.GetAuthUser(r)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	sessionID := r.PathValue("sessionID")
	session, err := s.sessionSrvc.GetByID(r.Context(), user, sessionID)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	id := r.PathValue("id")
	question, err := s.sessionQuestionSrvc.Update(
		r.Context(),
		session,
		id,
		req.Answer,
	)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusOK, question)
}
