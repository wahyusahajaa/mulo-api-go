package models

import "database/sql"

type UserVerified struct {
	Code      string
	ExpiredAt sql.NullTime
}
