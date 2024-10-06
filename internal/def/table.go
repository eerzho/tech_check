package def

type TableName string

const (
	TableUsers            TableName = "users"
	TableRefreshTokens    TableName = "refresh_tokens"
	TableRoles            TableName = "roles"
	TablePermissions      TableName = "permissions"
	TableCategories       TableName = "categories"
	TableQuestions        TableName = "questions"
	TableSessions         TableName = "sessions"
	TableSessionQuestions TableName = "session_questions"
)

func (tn TableName) String() string {
	return string(tn)
}
