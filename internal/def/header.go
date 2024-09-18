package def

type HeaderKey string

const (
	HeaderRequestID     HeaderKey = "X-Request-Id"
	HeaderContentType   HeaderKey = "Content-Type"
	HeaderForwardedFor  HeaderKey = "X-Forwarded-For"
	HeaderAuthorization HeaderKey = "Authorization"
)

func (hk HeaderKey) String() string {
	return string(hk)
}
