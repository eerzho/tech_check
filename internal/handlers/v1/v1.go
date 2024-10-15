package v1

import (
	"fmt"
	"net/http"
	"strings"
	"tech_check/internal/app"
	"tech_check/internal/handlers/v1/middlewares"
	"tech_check/internal/handlers/v1/requests"
	"tech_check/internal/handlers/v1/responses"
)

func New(mux *http.ServeMux, app *app.App) http.Handler {
	requests.InitParser()
	responses.InitBuilder(app.Cfg.IsDebug, app.Lg)

	authMwr := middlewares.NewAuth(app.Srvcs.Auth, app.Srvcs.User)
	permissionMwr := middlewares.NewPermission(app.Srvcs.User)

	newUser(mux, authMwr, permissionMwr, app.Srvcs.User)
	newAuth(mux, authMwr, app.Srvcs.Auth)
	newRole(mux, authMwr, permissionMwr, app.Srvcs.Role)
	newPermission(mux, authMwr, permissionMwr, app.Srvcs.Permission)
	newCategory(mux, authMwr, permissionMwr, app.Srvcs.Category)
	newQuestion(mux, authMwr, permissionMwr, app.Srvcs.Question)
	newSession(mux, authMwr, app.Srvcs.Session)
	newSessionQuestion(mux, authMwr, app.Srvcs.Session, app.Srvcs.SessionQuestion)

	reqIDMwr := middlewares.NewRequestId()
	reqLgMwr := middlewares.NewRequestLogger(app.Lg)

	return reqIDMwr.Mwr(reqLgMwr.Mwr(mux))
}

func Url(method, url string) string {
	return fmt.Sprintf("%s /api/v1/%s", method, strings.Trim(url, "/"))
}
