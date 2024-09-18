package v1

import (
	"fmt"
	"net/http"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

type user struct {
	rp       *request.Parser
	rb       *response.Builder
	userSrvc UserSrvc
}

func newUser(
	mux *http.ServeMux,
	prefix string,
	rp *request.Parser,
	rb *response.Builder,
	authMwr *mwr.Auth,
	userSrvc UserSrvc,
) {
	prefix += "/users"
	u := user{
		rp:       rp,
		rb:       rb,
		userSrvc: userSrvc,
	}

	mux.HandleFunc("GET "+prefix, authMwr.MwrFunc(u.list))
	mux.HandleFunc("POST "+prefix, u.create)
	mux.HandleFunc("GET "+prefix+"/{id}", authMwr.MwrFunc(u.show))
	mux.HandleFunc("PATCH "+prefix+"/{id}", authMwr.MwrFunc(u.update))
	mux.HandleFunc("DELETE "+prefix+"/{id}", authMwr.MwrFunc(u.delete))
	mux.HandleFunc("POST "+prefix+"/{id}/roles/{roleID}", authMwr.MwrFunc(u.addRole))
	mux.HandleFunc("DELETE "+prefix+"/{id}/roles/{roleID}", authMwr.MwrFunc(u.removeRole))
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
// @Success 200 {object} response.list{data=[]model.User,pagination=dto.Pagination}
func (u *user) list(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.list"

	search := u.rp.GetQuerySearch(r)
	users, pagination, err := u.userSrvc.List(
		r.Context(),
		search.Pagination.Page,
		search.Pagination.Count,
		search.Filters,
		search.Sorts,
	)
	if err != nil {
		u.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	u.rb.JsonList(w, r, users, pagination)
}

// @Summary registration
// @Tags users
// @Router /v1/users [post]
// @Accept json
// @Param body body request.UserCreate true "user create request"
// @Produce json
// @Success 201 {object} response.success{data=model.User}
func (u *user) create(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.create"

	var req request.UserCreate
	err := u.rp.ParseBody(r, &req)
	if err != nil {
		u.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	user, err := u.userSrvc.Create(
		r.Context(),
		req.Email,
		req.Name,
		req.Password,
	)
	if err != nil {
		u.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	u.rb.JsonSuccess(w, r, http.StatusCreated, user)
}

// @Summary get user by id
// @Tags users
// @Security BearerAuth
// @Router /v1/users/{id} [get]
// @Param id path string true "user id"
// @Produce json
// @Success 200 {object} response.success{data=model.User}
func (u *user) show(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.show"

	id := r.PathValue("id")
	user, err := u.userSrvc.GetByID(r.Context(), id)
	if err != nil {
		u.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	u.rb.JsonSuccess(w, r, http.StatusOK, user)
}

// @Summary update profile
// @Tags users
// @Security BearerAuth
// @Router /v1/users/{id} [patch]
// @Accept json
// @Param id path string true "user id"
// @Param body body request.UserUpdate true "user update request"
// @Produce json
// @Success 200 {object} response.success{data=model.User}
func (u *user) update(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.update"

	id := r.PathValue("id")
	var req request.UserUpdate

	err := u.rp.ParseBody(r, &req)
	if err != nil {
		u.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	user, err := u.userSrvc.Update(r.Context(), id, req.Name)
	if err != nil {
		u.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	u.rb.JsonSuccess(w, r, http.StatusOK, user)
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
		u.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	u.rb.JsonSuccess(w, r, http.StatusNoContent, nil)
}

// @Summary add role
// @Tags users
// @Security BearerAuth
// @Router /v1/users/{id}/roles/{roleID} [post]
// @Param id path string true "user id"
// @Param roleID path string true "role id"
// @Success 200 {object} response.success{data=model.User}
func (u *user) addRole(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.addRole"

	id := r.PathValue("id")
	roleID := r.PathValue("roleID")
	user, err := u.userSrvc.AddRole(r.Context(), id, roleID)
	if err != nil {
		u.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	u.rb.JsonSuccess(w, r, http.StatusOK, user)
}

// @Summary remove role
// @Tags users
// @Security BearerAuth
// @Router /v1/users/{id}/roles/{roleID} [delete]
// @Param id path string true "user id"
// @Param roleID path string true "role id"
// @Success 200 {object} response.success{data=model.User}
func (u *user) removeRole(w http.ResponseWriter, r *http.Request) {
	const op = "v1.user.removeRole"

	id := r.PathValue("id")
	roleID := r.PathValue("roleID")
	user, err := u.userSrvc.RemoveRole(r.Context(), id, roleID)
	if err != nil {
		u.rb.JsonFail(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	u.rb.JsonSuccess(w, r, http.StatusOK, user)
}
