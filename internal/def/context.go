package def

type ContextKey string

const (
	ContextAuthUser ContextKey = "auth_user"
)

func (ck ContextKey) String() string {
	return string(ck)
}
