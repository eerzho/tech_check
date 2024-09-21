package v1

import (
	"net/http"
	"tech_check/internal/app"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

func New(mux *http.ServeMux, app *app.App, prefix string) http.Handler {
	request.InitParser()
	response.InitBuilder(app.Cfg.IsDebug, app.Lg)

	authMwr := mwr.NewAuth(app.Srvcs.Auth, app.Srvcs.User)

	newUser(mux, prefix, authMwr, app.Srvcs.User)
	newAuth(mux, prefix, authMwr, app.Srvcs.Auth)
	newRole(mux, prefix, authMwr, app.Srvcs.Role)
	newPermission(mux, prefix, authMwr, app.Srvcs.Permission)
	newCategory(mux, prefix, authMwr, app.Srvcs.Category)
	newQuestion(mux, prefix, authMwr, app.Srvcs.Question)

	reqIDMwr := mwr.NewRequestId()
	reqLgMwr := mwr.NewRequestLogger(app.Lg)

	return reqIDMwr.Mwr(reqLgMwr.Mwr(mux))
}
