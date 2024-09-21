package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

type auth struct {
	authSrvc AuthSrvc
}

func newAuth(
	mux *http.ServeMux,
	prefix string,
	authMwr *mwr.Auth,
	authSrvc AuthSrvc,
) {
	prefix += "/auth"
	a := auth{
		authSrvc: authSrvc,
	}

	mux.HandleFunc("POST "+prefix, a.login)
	mux.HandleFunc("POST "+prefix+"/google", a.googleLogin)
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
	err := request.ParseBody(r, &req)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	token, err := a.authSrvc.Login(r.Context(), req.Email, req.Password, request.GetHeaderIP(r))
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusOK, token)
}

// @Summary google login
// @Tags auth
// @Router /v1/auth/google [post]
// @Accept json
// @Param body body request.GoogleLogin true "google login request"
// @Produce json
// @Success 200 {object} response.success{data=dto.Token}
func (a *auth) googleLogin(w http.ResponseWriter, r *http.Request) {
	const op = "v1.auth.googleLogin"

	var req request.GoogleLogin
	err := request.ParseBody(r, &req)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	token, err := a.authSrvc.GoogleLogin(r.Context(), req.TokenID, request.GetHeaderIP(r))
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusOK, token)
}

// @Summary get auth user
// @Tags auth
// @Security BearerAuth
// @Router /v1/auth [get]
// @Produce json
// @Success 200 {object} response.success{data=model.User}
func (a *auth) me(w http.ResponseWriter, r *http.Request) {
	const op = "v1.auth.me"

	user, err := request.GetAuthUser(r)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusOK, user)
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
	err := request.ParseBody(r, &req)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	token, err := a.authSrvc.Refresh(r.Context(), req.AToken, req.RToken, request.GetHeaderIP(r))
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusOK, token)
}
