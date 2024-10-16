package constants

import "errors"

var (
	ErrNotFound             = errors.New("resource not found")
	ErrAlreadyExists        = errors.New("resource already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidBody          = errors.New("request body cannot be empty")
	ErrAuthMissing          = errors.New("authorization header is missing")
	ErrInvalidAuthFormat    = errors.New("invalid authorization header format")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidClaimsType    = errors.New("invalid claims type")
	ErrATokenExpired        = errors.New("access token expired")
	ErrInvalidUserType      = errors.New("invalid user type")
	ErrCannotLogin          = errors.New("user cannot login")
	ErrTokensMismatch       = errors.New("refresh token and access token mismatch")
	ErrRTokenExpired        = errors.New("refresh token expired")
	ErrInvalidRToken        = errors.New("invalid refresh token")
	ErrInvalidGoogleData    = errors.New("invalid google data")
	ErrInvalidGradeValue    = errors.New("invalid grade value")
	ErrValidation           = errors.New("validation error")
	ErrAccessDenied         = errors.New("access denied")
	ErrUserHasActiveSession = errors.New("user has active session")
	ErrQuestionNotEnough    = errors.New("questions not enough")
	ErrSessionFinished      = errors.New("session finished")
)
