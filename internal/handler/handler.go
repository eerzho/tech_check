package handler

import (
	"net/http"
	_ "tech_check/docs"
	"tech_check/internal/app"
	v1 "tech_check/internal/handler/v1"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @Title tech_check http api
// @Version 1.0
// @BasePath /api
// @SecurityDefinitions.apikey BearerAuth
// @In header
// @Name Authorization
func New(app *app.App) http.Handler {
	mux := http.NewServeMux()

	// handler
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	handler := v1.New(mux, app, "/api/v1")

	return handler
}
