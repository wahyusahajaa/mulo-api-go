package models

import (
	"database/sql"
)

type User struct {
	Id              int
	Fullname        string
	Email           string
	Username        string
	Password        string
	Role            string
	Image           []byte
	EmailVerifiedAt sql.NullTime
}
