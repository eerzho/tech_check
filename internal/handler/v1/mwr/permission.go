package mwr

import (
	"fmt"
	"net/http"
	"tech_check/internal/def"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
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
	const op = "v1.mwr.Permission.MwrFunc"
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := request.GetAuthUser(r)
		if err != nil {
			response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}

		has, err := p.userSrvc.HasPermission(r.Context(), user, permissionSlug)
		if err != nil {
			response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}

		if !has {
			response.JsonFail(w, r, fmt.Errorf("%s: %w", op, def.ErrAccessDenied))
			return
		}

		next.ServeHTTP(w, r)
	}
}
