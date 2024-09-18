package v1

import (
	"net/http"
	"tech_check/internal/app"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

func New(mux *http.ServeMux, app *app.App, prefix string) http.Handler {
	rp := request.NewParser()
	rb := response.NewBuilder(app.Cfg.IsDebug, app.Lg)

	authMwr := mwr.NewAuth(rb, app.Srvcs.Auth, app.Srvcs.User)

	newUser(mux, prefix, rp, rb, authMwr, app.Srvcs.User)
	newAuth(mux, prefix, rp, rb, authMwr, app.Srvcs.Auth)
	newRole(mux, prefix, rp, rb, authMwr, app.Srvcs.Role)
	newPermission(mux, prefix, rp, rb, authMwr, app.Srvcs.Permission)
	newCategory(mux, prefix, rp, rb, authMwr, app.Srvcs.Category)
	newQuestion(mux, prefix, rp, rb, authMwr, app.Srvcs.Question)

	reqIDMwr := mwr.NewRequestId()
	reqLgMwr := mwr.NewRequestLogger(app.Lg)

	return reqIDMwr.Mwr(reqLgMwr.Mwr(mux))
}
