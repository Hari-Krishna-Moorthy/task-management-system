package utils

import "time"

const (
	NumbeZero = 0
	Number14  = 14
	Number72  = 72
)

const (
	JWT_TOKEN_EXPIRY   = Number72 * time.Hour
	JWT_DEFAULT_SECRET = "default-secret-key"
)

const (
	Email           = "email"
	Username        = "username"
	CookieKeyToken  = "token"
	DeletedAt       = "deleted_at"
	MongoDBFilterOr = "$or"
)

// color codes
const (
	RedColor   = "\033[31m"
	GreenColor = "\033[32m"
	BlueColor  = "\033[34m"
	ResetColor = "\033[0m"
)
