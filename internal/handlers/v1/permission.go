package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handlers/v1/middlewares"
	"tech_check/internal/handlers/v1/requests"
	"tech_check/internal/handlers/v1/responses"
)

type permission struct {
	permissionService PermissionService
}

func newPermission(
	mux *http.ServeMux,
	authMwr *middlewares.Auth,
	permissionMwr *middlewares.Permission,
	permissionService PermissionService,
) {
	p := permission{
		permissionService: permissionService,
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
// @Success 200 {object} responses.list{data=[]models.Permission,pagination=dto.Pagination}
func (p *permission) list(w http.ResponseWriter, r *http.Request) {
	const op = "v1.permission.list"

	search := requests.GetQuerySearch(r)
	permissions, pagination, err := p.permissionService.List(
		r.Context(),
		search.Pagination.Page,
		search.Pagination.Count,
		search.Filters,
		search.Sorts,
	)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonList(w, r, permissions, pagination)
}

// @Summary get permission by id
// @Tags permissions
// @Security BearerAuth
// @Router /v1/permissions/{id} [get]
// @Param id path string true "permission id"
// @Produce json
// @Success 200 {object} responses.success{data=models.Permission}
func (p *permission) show(w http.ResponseWriter, r *http.Request) {
	const op = "v1.permission.show"

	id := r.PathValue("id")
	permission, err := p.permissionService.GetByID(r.Context(), id)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, permission)
}
