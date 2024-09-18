package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"tech_check/internal/def"
)

type (
	Builder struct {
		isDebug         bool
		lg              *slog.Logger
		strangeCaseJson string
	}

	fail struct {
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message"`
	}

	success struct {
		Data interface{} `json:"data,omitempty"`
	}

	list struct {
		Data       interface{} `json:"data"`
		Pagination interface{} `json:"pagination,omitempty"`
	}
)

func NewBuilder(isDebug bool, lg *slog.Logger) *Builder {
	return &Builder{
		isDebug:         isDebug,
		lg:              lg,
		strangeCaseJson: `{"message": "` + http.StatusText(http.StatusInternalServerError) + `"}`,
	}
}

func (b *Builder) JsonFail(w http.ResponseWriter, r *http.Request, err error) {
	code := http.StatusInternalServerError
	msg := b.originalErr(err).Error()

	if errors.Is(err, def.ErrNotFound) {
		code = http.StatusNotFound
	} else if errors.Is(err, def.ErrAlreadyExists) ||
		errors.Is(err, def.ErrInvalidBody) {
		code = http.StatusBadRequest
	} else if errors.Is(err, def.ErrInvalidCredentials) ||
		errors.Is(err, def.ErrAuthMissing) ||
		errors.Is(err, def.ErrInvalidAuthFormat) ||
		errors.Is(err, def.ErrInvalidSigningMethod) ||
		errors.Is(err, def.ErrInvalidClaimsType) ||
		errors.Is(err, def.ErrATokenExpired) ||
		errors.Is(err, def.ErrInvalidUserType) ||
		errors.Is(err, def.ErrTokensMismatch) ||
		errors.Is(err, def.ErrRTokenExpired) ||
		errors.Is(err, def.ErrInvalidRToken) {
		code = http.StatusUnauthorized
	} else if errors.Is(err, def.ErrCannotLogin) {
		code = http.StatusForbidden
	}

	if !b.isDebug && code == http.StatusInternalServerError {
		msg = http.StatusText(code)
	}

	f := fail{Message: msg}

	b.logFail(r, code, err)
	b.Json(w, code, &f)
}

func (b *Builder) JsonSuccess(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	s := success{Data: data}

	b.logSuccess(r, code)
	b.Json(w, code, &s)
}

func (b *Builder) JsonList(w http.ResponseWriter, r *http.Request, data, pagination interface{}) {
	l := list{
		Data:       data,
		Pagination: pagination,
	}

	b.logSuccess(r, http.StatusOK)
	b.Json(w, http.StatusOK, &l)
}

func (b *Builder) Json(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set(def.HeaderContentType.String(), "application/json")

	jsonBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, b.strangeCaseJson, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(jsonBody)
}

func (b *Builder) originalErr(err error) error {
	unwrappedErr := errors.Unwrap(err)
	if unwrappedErr == nil {
		return err
	}
	return b.originalErr(unwrappedErr)
}

func (b *Builder) logFail(r *http.Request, code int, err error) {
	lg := b.lg.With(
		slog.Any(def.HeaderRequestID.String(), r.Context().Value(def.HeaderRequestID)),
		slog.String("method", r.Method),
		slog.String("url", r.URL.Path),
		slog.Int("status", code),
	)

	if code >= http.StatusInternalServerError {
		lg.Error(err.Error())
	} else {
		lg.Debug(err.Error())
	}
}

func (b *Builder) logSuccess(r *http.Request, code int) {
	lg := b.lg.With(
		slog.Any(def.HeaderRequestID.String(), r.Context().Value(def.HeaderRequestID)),
		slog.String("method", r.Method),
		slog.String("url", r.URL.Path),
		slog.Int("status", code),
	)

	lg.Debug("success")
}
