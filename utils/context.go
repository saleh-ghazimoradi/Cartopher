package utils

type ContextKey string

const (
	UserIdKey     ContextKey = "user_id"
	UserEmailKey  ContextKey = "user_email"
	UserRoleKey   ContextKey = "user_role"
	GinContextKey ContextKey = "gin_context"
)
