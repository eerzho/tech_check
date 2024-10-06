package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

type session struct {
	sessionSrvc SessionSrvc
}

func newSession(
	mux *http.ServeMux,
	authMwr *mwr.Auth,
	sessionSrvc SessionSrvc,
) {
	s := session{
		sessionSrvc: sessionSrvc,
	}

	mux.HandleFunc(
		Url(http.MethodPost, "/sessions"),
		authMwr.MwrFunc(s.create),
	)
}

// @Summary start test session
// @Tags sessions
// @Security BearerAuth
// @Router /v1/sessions [post]
// @Param body body request.SessionCreate true "session create request"
// @Produce json
// @Success 201 {object} response.success{data=model.Session}
func (s *session) create(w http.ResponseWriter, r *http.Request) {
	const op = "v1.session.create"

	user, err := request.GetAuthUser(r)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	var req request.SessionCreate
	err = request.ParseBody(r, &req)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	session, err := s.sessionSrvc.Create(
		r.Context(),
		user,
		req.CategoryID,
		req.Grade,
	)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusCreated, session)
}
