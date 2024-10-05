package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

type permission struct {
	permissionSrvc PermissionSrvc
}

func newPermission(
	mux *http.ServeMux,
	authMwr *mwr.Auth,
	permissionMwr *mwr.Permission,
	permissionSrvc PermissionSrvc,
) {
	p := permission{
		permissionSrvc: permissionSrvc,
	}

	mux.HandleFunc(
		Url(http.MethodGet, "/permissions"),
		authMwr.MwrFunc(permissionMwr.MwrFunc(p.list, "permission-read")),
	)

	mux.HandleFunc(
		Url(http.MethodGet, "/permissions/{id}"),
		authMwr.MwrFunc(permissionMwr.MwrFunc(p.show, "permission-read")),
	)
}

// @Summary permissions list
// @Tags permissions
// @Security BearerAuth
// @Router /v1/permissions [get]
// @Param pagination[page] query int false "page"
// @Param pagination[count] query int false "count"
// @Param sorts[created_at] query string false "created_at" Enums(asc, desc)
// @Param sorts[updated_at] query string false "updated_at" Enums(asc, desc)
// @Param sorts[name] query string false "name" Enums(asc, desc)
// @Param sorts[slug] query string false "slug" Enums(asc, desc)
// @Param filters[name] query string false "name"
// @Param filters[slug] query string false "slug"
// @Produce json
// @Success 200 {object} response.list{data=[]model.Permission,pagination=dto.Pagination}
func (p *permission) list(w http.ResponseWriter, r *http.Request) {
	const op = "v1.permission.list"

	search := request.GetQuerySearch(r)
	permissions, pagination, err := p.permissionSrvc.List(
		r.Context(),
		search.Pagination.Page,
		search.Pagination.Count,
		search.Filters,
		search.Sorts,
	)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonList(w, r, permissions, pagination)
}

// @Summary get permission by id
// @Tags permissions
// @Security BearerAuth
// @Router /v1/permissions/{id} [get]
// @Param id path string true "permission id"
// @Produce json
// @Success 200 {object} response.success{data=model.Permission}
func (p *permission) show(w http.ResponseWriter, r *http.Request) {
	const op = "v1.permission.show"

	id := r.PathValue("id")
	permission, err := p.permissionSrvc.GetByID(r.Context(), id)
	if err != nil {
		response.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	response.JsonSuccess(w, r, http.StatusOK, permission)
}
