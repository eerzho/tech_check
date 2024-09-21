package response

import (
	"log/slog"
	"net/http"
)

var defaultBuilder *Builder

func InitBuilder(isDebug bool, lg *slog.Logger) {
	defaultBuilder = NewBuilder(isDebug, lg)
}

func JsonFail(w http.ResponseWriter, r *http.Request, err error) {
	defaultBuilder.JsonFail(w, r, err)
}

func JsonSuccess(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	defaultBuilder.JsonSuccess(w, r, code, data)
}

func JsonList(w http.ResponseWriter, r *http.Request, data, pagination interface{}) {
	defaultBuilder.JsonList(w, r, data, pagination)
}

func Json(w http.ResponseWriter, code int, body interface{}) {
	defaultBuilder.Json(w, code, body)
}
