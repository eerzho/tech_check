package middlewares

import (
	"fmt"
	"net/http"
	"tech_check/internal/constants"
	"tech_check/internal/handlers/v1/requests"
	"tech_check/internal/handlers/v1/responses"
)

type Permission struct {
	userSrvc UserSrvc
}

func NewPermission(userSrvc UserSrvc) *Permission {
	return &Permission{
		userSrvc: userSrvc,
	}
}

func (p *Permission) MwrFunc(next http.HandlerFunc, permissionSlug string) http.HandlerFunc {
	const op = "v1.middlewares.Permission.MwrFunc"
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.GetAuthUser(r)
		if err != nil {
			responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}

		has, err := p.userSrvc.HasPermission(r.Context(), user, permissionSlug)
		if err != nil {
			responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}

		if !has {
			responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, constants.ErrAccessDenied))
			return
		}

		next.ServeHTTP(w, r)
	}
}
