package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handlers/v1/middlewares"
	"tech_check/internal/handlers/v1/requests"
	"tech_check/internal/handlers/v1/responses"
)

type user struct {
	userSrvc UserSrvc
}

func newUser(
	mux *http.ServeMux,
	authMwr *middlewares.Auth,
	permissionMwr *middlewares.Permission,
	userSrvc UserSrvc,
) {
	u := user{
		userSrvc: userSrvc,
	}

	mux.HandleFunc(
		Url(http.MethodGet, "/users"),
		authMwr.MwrFunc(permissionMwr.MwrFunc(u.list, "user-read")),
	)

	mux.HandleFunc(
		Url(http.MethodPost, "/users"),
		u.create,
	)

	mux.HandleFunc(
		Url(http.MethodGet, "/users/{id}"),
		authMwr.MwrFunc(permissionMwr.MwrFunc(u.show, "user-read")),
	)

	mux.HandleFunc(
		Url(http.MethodPatch, "/users/{id}"),
		authMwr.MwrFunc(permissionMwr.MwrFunc(u.update, "user-edit")),
	)

	mux.HandleFunc(
		Url(http.MethodDelete, "/users/{id}"),
		authMwr.MwrFunc(permissionMwr.MwrFunc(u.delete, "user-delete")),
	)

	mux.HandleFunc(
		Url(http.MethodPost, "/users/{id}/roles/{roleID}"),
		authMwr.MwrFunc(permissionMwr.MwrFunc(u.addRole, "user-edit")),
	)

	mux.HandleFunc(
		Url(http.MethodDelete, "/users/{id}/roles/{roleID}"),
		authMwr.MwrFunc(permissionMwr.MwrFunc(u.removeRole, "user-edit")),
	)
}

// @Summary users list
// @Tags users
// @Security BearerAuth
// @Router /v1/users [get]
// @Param pagination[page] query int false "page"
// @Param pagination[count] query int false "count"
// @Param sorts[created_at] query string false "created_at" Enums(asc, desc)
// @Param sorts[updated_at] query string false "updated_at" Enums(asc, desc)
// @Param filters[name] query string false "name"
// @Param filters[email] query string false "email"
// @Produce json
// @Success 200 {object} responses.list{data=[]models.User,pagination=dto.Pagination}
func (u *user) list(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.list"

	search := requests.GetQuerySearch(r)
	users, pagination, err := u.userSrvc.List(
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

	responses.JsonList(w, r, users, pagination)
}

// @Summary registration
// @Tags users
// @Router /v1/users [post]
// @Accept json
// @Param body body requests.UserCreate true "user create request"
// @Produce json
// @Success 201 {object} responses.success{data=models.User}
func (u *user) create(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.create"

	var req requests.UserCreate
	err := requests.ParseBody(r, &req)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	user, err := u.userSrvc.Create(
		r.Context(),
		req.Email,
		req.Name,
		req.Password,
	)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusCreated, user)
}

// @Summary get user by id
// @Tags users
// @Security BearerAuth
// @Router /v1/users/{id} [get]
// @Param id path string true "user id"
// @Produce json
// @Success 200 {object} responses.success{data=models.User}
func (u *user) show(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.show"

	id := r.PathValue("id")
	user, err := u.userSrvc.GetByID(r.Context(), id)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, user)
}

// @Summary update profile
// @Tags users
// @Security BearerAuth
// @Router /v1/users/{id} [patch]
// @Accept json
// @Param id path string true "user id"
// @Param body body requests.UserUpdate true "user update request"
// @Produce json
// @Success 200 {object} responses.success{data=models.User}
func (u *user) update(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.update"

	id := r.PathValue("id")
	var req requests.UserUpdate

	err := requests.ParseBody(r, &req)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	user, err := u.userSrvc.Update(r.Context(), id, req.Name)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, user)
}

// @Summary delete user by id
// @Tags users
// @Security BearerAuth
// @Router /v1/users/{id} [delete]
// @Param id path string true "user id"
// @Success 204
func (u *user) delete(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.delete"

	id := r.PathValue("id")

	err := u.userSrvc.Delete(r.Context(), id)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusNoContent, nil)
}

// @Summary add role
// @Tags users
// @Security BearerAuth
// @Router /v1/users/{id}/roles/{roleID} [post]
// @Param id path string true "user id"
// @Param roleID path string true "role id"
// @Success 200 {object} responses.success{data=models.User}
func (u *user) addRole(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.addRole"

	id := r.PathValue("id")
	roleID := r.PathValue("roleID")
	user, err := u.userSrvc.AddRole(r.Context(), id, roleID)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, user)
}

// @Summary remove role
// @Tags users
// @Security BearerAuth
// @Router /v1/users/{id}/roles/{roleID} [delete]
// @Param id path string true "user id"
// @Param roleID path string true "role id"
// @Success 200 {object} responses.success{data=models.User}
func (u *user) removeRole(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.removeRole"

	id := r.PathValue("id")
	roleID := r.PathValue("roleID")
	user, err := u.userSrvc.RemoveRole(r.Context(), id, roleID)
	if err != nil {
		responses.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	responses.JsonSuccess(w, r, http.StatusOK, user)
}
