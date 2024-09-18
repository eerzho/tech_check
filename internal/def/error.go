package def

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
)
