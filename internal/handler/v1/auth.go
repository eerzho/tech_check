package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

type auth struct {
	rp       *request.Parser
	rb       *response.Builder
	authSrvc AuthSrvc
}

func newAuth(
	mux *http.ServeMux,
	prefix string,
	rp *request.Parser,
	rb *response.Builder,
	authMwr *mwr.Auth,
	authSrvc AuthSrvc,
) {
	prefix += "/auth"
	a := auth{
		rp:       rp,
		rb:       rb,
		authSrvc: authSrvc,
	}

	mux.HandleFunc("POST "+prefix, a.login)
	mux.HandleFunc("GET "+prefix, authMwr.MwrFunc(a.me))
	mux.HandleFunc("POST "+prefix+"/refresh", a.refresh)
}

// @Summary login
// @Tags auth
// @Router /v1/auth [post]
// @Accept json
// @Param body body request.Login true "login request"
// @Produce json
// @Success 200 {object} response.success{data=dto.Token}
func (a *auth) login(w http.ResponseWriter, r *http.Request) {
	const op = "v1.auth.login"

	var req request.Login
	err := a.rp.ParseBody(r, &req)
	if err != nil {
		a.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	token, err := a.authSrvc.Login(r.Context(), req.Email, req.Password, a.rp.GetHeaderIP(r))
	if err != nil {
		a.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	a.rb.JsonSuccess(w, r, http.StatusOK, token)
}

// @Summary get auth user
// @Tags auth
// @Security BearerAuth
// @Router /v1/auth [get]
// @Produce json
// @Success 200 {object} response.success{data=model.User}
func (a *auth) me(w http.ResponseWriter, r *http.Request) {
	const op = "v1.auth.me"

	user, err := a.rp.GetAuthUser(r)
	if err != nil {
		a.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	a.rb.JsonSuccess(w, r, http.StatusOK, user)
}

// @Summary refresh token
// @Tags auth
// @Router /v1/auth/refresh [post]
// @Accept json
// @Param body body request.Refresh true "token refresh request"
// @Produce json
// @Success 200 {object} response.success{data=dto.Token}
func (a *auth) refresh(w http.ResponseWriter, r *http.Request) {
	const op = "v1.auth.refresh"

	var req request.Refresh
	err := a.rp.ParseBody(r, &req)
	if err != nil {
		a.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	token, err := a.authSrvc.Refresh(r.Context(), req.AToken, req.RToken, a.rp.GetHeaderIP(r))
	if err != nil {
		a.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	a.rb.JsonSuccess(w, r, http.StatusOK, token)
}
