package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handlers/v1/middlewares"
	"tech_check/internal/handlers/v1/requests"
	"tech_check/internal/handlers/v1/responses"
)

type sessionQuestion struct {
	sessionQuestionService SessionQuestionService
}

func newSessionQuestion(
	mux *http.ServeMux,
	authMwr *middlewares.Auth,
	sessionQuestionService SessionQuestionService,
) {
	s := sessionQuestion{
		sessionQuestionService: sessionQuestionService,
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
// @Success 200 {object} responses.success{data=[]models.SessionQuestion}
func (s *sessionQuestion) list(w http.ResponseWriter, r *http.Request) {
	const op = "v1.sessionQuestion.list"

	user, err := requests.GetAuthUser(r)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	sessionID := r.PathValue("sessionID")
	questions, err := s.sessionQuestionService.List(r.Context(), user, sessionID)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, questions)
}

// @Summary get session question by id
// @Tags sessionQuestions
// @Security BearerAuth
// @Router /v1/sessions/{sessionID}/questions/{id} [get]
// @Param sessionID path string true "session id"
// @Param id path string true "session question id"
// @Produce json
// @Success 200 {object} responses.success{data=models.SessionQuestion}
func (s *sessionQuestion) show(w http.ResponseWriter, r *http.Request) {
	const op = "v1.sessionQuestion.show"

	user, err := requests.GetAuthUser(r)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	id := r.PathValue("id")
	sessionID := r.PathValue("sessionID")
	question, err := s.sessionQuestionService.GetByID(r.Context(), user, sessionID, id)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, question)
}

// @Summary answer the question
// @Tags sessionQuestions
// @Security BearerAuth
// @Router /v1/sessions/{sessionID}/questions/{id} [patch]
// @Accept json
// @Param sessionID path string true "session id"
// @Param id path string true "session question id"
// @Param body body requests.SessionQuestionUpdate true "answer the question request"
// @Produce json
// @Success 200 {object} responses.success{data=models.SessionQuestion}
func (s *sessionQuestion) update(w http.ResponseWriter, r *http.Request) {
	const op = "v1.sessionQuestion.update"

	var req requests.SessionQuestionUpdate
	err := requests.ParseBody(r, &req)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	user, err := requests.GetAuthUser(r)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	id := r.PathValue("id")
	sessionID := r.PathValue("sessionID")
	question, err := s.sessionQuestionService.Update(
		r.Context(),
		user,
		sessionID,
		id,
		req.Answer,
	)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, question)
}
