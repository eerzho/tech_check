package responses

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"tech_check/internal/constants"

	"github.com/go-playground/validator/v10"
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
	var code int
	var msg string
	var data interface{}

	ve, ok := b.originalErr(err).(validator.ValidationErrors)
	if ok {
		code = http.StatusBadRequest
		msg = constants.ErrValidation.Error()
		data = b.getData(ve)
	} else {
		code = b.getCode(err)
		msg = b.getMsg(code, err)
		data = nil
	}

	f := fail{Message: msg, Data: data}

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
	w.Header().Set(constants.HeaderContentType.String(), "application/json")

	jsonBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, b.strangeCaseJson, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(jsonBody)
}

func (b *Builder) getCode(err error) int {
	code := http.StatusInternalServerError

	if errors.Is(err, constants.ErrNotFound) {
		code = http.StatusNotFound
	} else if errors.Is(err, constants.ErrAlreadyExists) ||
		errors.Is(err, constants.ErrInvalidBody) ||
		errors.Is(err, constants.ErrUserHasActiveSession) ||
		errors.Is(err, constants.ErrQuestionNotEnough) ||
		errors.Is(err, constants.ErrSessionFinished) {
		code = http.StatusBadRequest
	} else if errors.Is(err, constants.ErrInvalidCredentials) ||
		errors.Is(err, constants.ErrAuthMissing) ||
		errors.Is(err, constants.ErrInvalidAuthFormat) ||
		errors.Is(err, constants.ErrInvalidSigningMethod) ||
		errors.Is(err, constants.ErrInvalidClaimsType) ||
		errors.Is(err, constants.ErrATokenExpired) ||
		errors.Is(err, constants.ErrInvalidUserType) ||
		errors.Is(err, constants.ErrTokensMismatch) ||
		errors.Is(err, constants.ErrRTokenExpired) ||
		errors.Is(err, constants.ErrInvalidRToken) {
		code = http.StatusUnauthorized
	} else if errors.Is(err, constants.ErrCannotLogin) ||
		errors.Is(err, constants.ErrAccessDenied) {
		code = http.StatusForbidden
	}

	return code
}

func (b *Builder) getMsg(code int, err error) string {
	if b.isDebug {
		return b.originalErr(err).Error()
	}

	if code >= http.StatusInternalServerError {
		return http.StatusText(code)
	}

	return b.originalErr(err).Error()
}

func (b *Builder) getData(ve validator.ValidationErrors) map[string]string {
	data := make(map[string]string, len(ve))

	for _, err := range ve {
		data[err.Field()] = fmt.Sprintf("validation failed on the '%s' tag", err.Tag())
	}

	return data
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
		slog.Any(constants.HeaderRequestID.String(), r.Context().Value(constants.HeaderRequestID)),
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
		slog.Any(constants.HeaderRequestID.String(), r.Context().Value(constants.HeaderRequestID)),
		slog.String("method", r.Method),
		slog.String("url", r.URL.Path),
		slog.Int("status", code),
	)

	lg.Debug("success")
}
