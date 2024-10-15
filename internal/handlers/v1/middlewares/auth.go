package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"tech_check/internal/constants"
	"tech_check/internal/handlers/v1/responses"
	"time"
)

type Auth struct {
	authSrvc AuthSrvc
	userSrvc UserSrvc
}

func NewAuth(
	authSrvc AuthSrvc,
	userSrvc UserSrvc,
) *Auth {
	return &Auth{
		authSrvc: authSrvc,
		userSrvc: userSrvc,
	}
}

func (a *Auth) MwrFunc(next http.HandlerFunc) http.HandlerFunc {
	const op = "v1.middlewares.Auth.MwrFunc"
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(constants.HeaderAuthorization.String())
		if authHeader == "" {
			responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, constants.ErrAuthMissing))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, constants.ErrInvalidAuthFormat))
			return
		}

		claims, err := a.authSrvc.DecodeAToken(r.Context(), parts[1])
		if err != nil {
			responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, constants.ErrATokenExpired))
			return
		}

		user, err := a.userSrvc.GetByID(r.Context(), claims.UserID)
		if err != nil {
			if errors.Is(err, constants.ErrNotFound) {
				responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, constants.ErrCannotLogin))
				return
			}
			responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}

		ctx := context.WithValue(r.Context(), constants.ContextAuthUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
