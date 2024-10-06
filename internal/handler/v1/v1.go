package v1

import (
	"fmt"
	"net/http"
	"strings"
	"tech_check/internal/app"
	"tech_check/internal/handler/v1/mwr"
	"tech_check/internal/handler/v1/request"
	"tech_check/internal/handler/v1/response"
)

func New(mux *http.ServeMux, app *app.App) http.Handler {
	request.InitParser()
	response.InitBuilder(app.Cfg.IsDebug, app.Lg)

	authMwr := mwr.NewAuth(app.Srvcs.Auth, app.Srvcs.User)
	permissionMwr := mwr.NewPermission(app.Srvcs.User)

	newUser(mux, authMwr, permissionMwr, app.Srvcs.User)
	newAuth(mux, authMwr, app.Srvcs.Auth)
	newRole(mux, authMwr, permissionMwr, app.Srvcs.Role)
	newPermission(mux, authMwr, permissionMwr, app.Srvcs.Permission)
	newCategory(mux, authMwr, permissionMwr, app.Srvcs.Category)
	newQuestion(mux, authMwr, permissionMwr, app.Srvcs.Question)
	newSession(mux, authMwr, app.Srvcs.Session)

	reqIDMwr := mwr.NewRequestId()
	reqLgMwr := mwr.NewRequestLogger(app.Lg)

	return reqIDMwr.Mwr(reqLgMwr.Mwr(mux))
}

func Url(method, url string) string {
	return fmt.Sprintf("%s /api/v1/%s", method, strings.Trim(url, "/"))
}
