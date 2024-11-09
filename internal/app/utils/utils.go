package utils

import (
	"time"
)

const (
	NumbeZero = 0
	Number14  = 14
	Number72  = 72

	TimeLayout  = "2006-01-02"
	EmptyString = ""
)

const (
	JWT_TOKEN_EXPIRY   = Number72 * time.Hour
	JWT_DEFAULT_SECRET = "default-secret-key"
)

const (
	IdKey = "_id"
	UserIDKey = "user_id"
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

var DeleteAtZeroTime = time.Time{}.UTC()
