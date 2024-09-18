package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

type role struct {
	rp       *request.Parser
	rb       *response.Builder
	roleSrvc RoleSrvc
}

func newRole(
	mux *http.ServeMux,
	prefix string,
	rp *request.Parser,
	rb *response.Builder,
	authMwr *mwr.Auth,
	roleSrvc RoleSrvc,
) {
	prefix += "/roles"
	re := role{
		rp:       rp,
		rb:       rb,
		roleSrvc: roleSrvc,
	}

	mux.HandleFunc("GET "+prefix, authMwr.MwrFunc(re.list))
	mux.HandleFunc("POST "+prefix, authMwr.MwrFunc(re.create))
	mux.HandleFunc("GET "+prefix+"/{id}", authMwr.MwrFunc(re.show))
	mux.HandleFunc("PATCH "+prefix+"/{id}", authMwr.MwrFunc(re.update))
	mux.HandleFunc("DELETE "+prefix+"/{id}", authMwr.MwrFunc(re.delete))
	mux.HandleFunc("POST "+prefix+"/{id}/permissions/{permissionID}", authMwr.MwrFunc(re.addPermission))
	mux.HandleFunc("DELETE "+prefix+"/{id}/permissions/{permissionID}", authMwr.MwrFunc(re.removePermission))
}

// @Summary roles list
// @Tags roles
// @Security BearerAuth
// @Router /v1/roles [get]
// @Param pagination[page] query int false "page"
// @Param pagination[count] query int false "count"
// @Param sorts[created_at] query string false "created_at" Enums(asc, desc)
// @Param sorts[updated_at] query string false "updated_at" Enums(asc, desc)
// @Param sorts[name] query string false "name" Enums(asc, desc)
// @Param sorts[slug] query string false "slug" Enums(asc, desc)
// @Param filters[name] query string false "name"
// @Param filters[slug] query string false "slug"
// @Produce json
// @Success 200 {object} response.list{data=[]model.Role,pagination=dto.Pagination}
func (re *role) list(w http.ResponseWriter, r *http.Request) {
	const op = "v1.role.list"

	search := re.rp.GetQuerySearch(r)
	roles, pagination, err := re.roleSrvc.List(
		r.Context(),
		search.Pagination.Page,
		search.Pagination.Count,
		search.Filters,
		search.Sorts,
	)
	if err != nil {
		re.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	re.rb.JsonList(w, r, roles, pagination)
}

// @Summary create role
// @Tags roles
// @Security BearerAuth
// @Router /v1/roles [post]
// @Accept json
// @Param body body request.RoleCreate true "role create request"
// @Produce json
// @Success 201 {object} response.success{data=model.Role}
func (re *role) create(w http.ResponseWriter, r *http.Request) {
	const op = "v1.role.create"

	var req request.RoleCreate
	err := re.rp.ParseBody(r, &req)
	if err != nil {
		re.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	role, err := re.roleSrvc.Create(r.Context(), req.Name)
	if err != nil {
		re.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	re.rb.JsonSuccess(w, r, http.StatusCreated, role)
}

// @Summary get role by id
// @Tags roles
// @Security BearerAuth
// @Router /v1/roles/{id} [get]
// @Param id path string true "role id"
// @Produce json
// @Success 200 {object} response.success{data=model.Role}
func (re *role) show(w http.ResponseWriter, r *http.Request) {
	const op = "v1.role.show"

	id := r.PathValue("id")
	role, err := re.roleSrvc.GetByID(r.Context(), id)
	if err != nil {
		re.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	re.rb.JsonSuccess(w, r, http.StatusOK, role)
}

// @Summary update role by id
// @Tags roles
// @Security BearerAuth
// @Router /v1/roles/{id} [patch]
// @Accept json
// @Param id path string true "role id"
// @Param body body request.RoleUpdate true "role update request"
// @Produce json
// @Success 200 {object} response.success{data=model.Role}
func (re *role) update(w http.ResponseWriter, r *http.Request) {
	const op = "v1.role.update"

	id := r.PathValue("id")
	var req request.RoleUpdate

	err := re.rp.ParseBody(r, &req)
	if err != nil {
		re.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	role, err := re.roleSrvc.Update(r.Context(), id, req.Name)
	if err != nil {
		re.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	re.rb.JsonSuccess(w, r, http.StatusOK, role)
}

// @Summary delete role by id
// @Tags roles
// @Security BearerAuth
// @Router /v1/roles/{id} [delete]
// @Param id path string true "role id"
// @Success 204
func (re *role) delete(w http.ResponseWriter, r *http.Request) {
	const op = "v1.role.delete"

	id := r.PathValue("id")
	err := re.roleSrvc.Delete(r.Context(), id)
	if err != nil {
		re.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	re.rb.JsonSuccess(w, r, http.StatusNoContent, nil)
}

// @Summary add permission
// @Tags roles
// @Security BearerAuth
// @Router /v1/roles/{id}/permissions/{permissionID} [post]
// @Param id path string true "role id"
// @Param permissionID path string true "permission id"
// @Success 200 {object} response.success{data=model.Role}
func (re *role) addPermission(w http.ResponseWriter, r *http.Request) {
	const op = "v1.role.addPermission"

	id := r.PathValue("id")
	permissionID := r.PathValue("permissionID")
	role, err := re.roleSrvc.AddPermission(r.Context(), id, permissionID)
	if err != nil {
		re.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	re.rb.JsonSuccess(w, r, http.StatusOK, role)
}

// @Summary remove permission
// @Tags roles
// @Security BearerAuth
// @Router /v1/roles/{id}/permissions/{permissionID} [delete]
// @Param id path string true "role id"
// @Param permissionID path string true "permission id"
// @Success 200 {object} response.success{data=model.Role}
func (re *role) removePermission(w http.ResponseWriter, r *http.Request) {
	const op = "v1.role.removePermission"

	id := r.PathValue("id")
	permissionID := r.PathValue("permissionID")
	role, err := re.roleSrvc.RemovePermission(r.Context(), id, permissionID)
	if err != nil {
		re.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	re.rb.JsonSuccess(w, r, http.StatusOK, role)
}
