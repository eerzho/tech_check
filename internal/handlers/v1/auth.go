package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handlers/v1/middlewares"
	"tech_check/internal/handlers/v1/requests"
	"tech_check/internal/handlers/v1/responses"
)

type auth struct {
	authService AuthService
}

func newAuth(
	mux *http.ServeMux,
	authMwr *middlewares.Auth,
	authService AuthService,
) {
	a := auth{
		authService: authService,
	}

	mux.HandleFunc(Url(http.MethodPost, "/auth"), a.login)
	mux.HandleFunc(Url(http.MethodPost, "/auth/google"), a.googleLogin)
	mux.HandleFunc(Url(http.MethodGet, "/auth"), authMwr.MwrFunc(a.me))
	mux.HandleFunc(Url(http.MethodPost, "/auth/refresh"), a.refresh)
}

// @Summary login
// @Tags auth
// @Router /v1/auth [post]
// @Accept json
// @Param body body requests.Login true "login request"
// @Produce json
// @Success 200 {object} responses.success{data=dto.Token}
func (a *auth) login(w http.ResponseWriter, r *http.Request) {
	const op = "v1.auth.login"

	var req requests.Login
	err := requests.ParseBody(r, &req)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	token, err := a.authService.Login(r.Context(), req.Email, req.Password, requests.GetHeaderIP(r))
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, token)
}

// @Summary google login
// @Tags auth
// @Router /v1/auth/google [post]
// @Accept json
// @Param body body requests.GoogleLogin true "google login request"
// @Produce json
// @Success 200 {object} responses.success{data=dto.Token}
func (a *auth) googleLogin(w http.ResponseWriter, r *http.Request) {
	const op = "v1.auth.googleLogin"

	var req requests.GoogleLogin
	err := requests.ParseBody(r, &req)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	token, err := a.authService.GoogleLogin(r.Context(), req.TokenID, requests.GetHeaderIP(r))
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, token)
}

// @Summary get auth user
// @Tags auth
// @Security BearerAuth
// @Router /v1/auth [get]
// @Produce json
// @Success 200 {object} responses.success{data=models.User}
func (a *auth) me(w http.ResponseWriter, r *http.Request) {
	const op = "v1.auth.me"

	user, err := requests.GetAuthUser(r)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, user)
}

// @Summary refresh token
// @Tags auth
// @Router /v1/auth/refresh [post]
// @Accept json
// @Param body body requests.Refresh true "token refresh request"
// @Produce json
// @Success 200 {object} responses.success{data=dto.Token}
func (a *auth) refresh(w http.ResponseWriter, r *http.Request) {
	const op = "v1.auth.refresh"

	var req requests.Refresh
	err := requests.ParseBody(r, &req)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	token, err := a.authService.Refresh(r.Context(), req.AToken, req.RToken, requests.GetHeaderIP(r))
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, token)
}
