package models

import (
	"database/sql"
)

type User struct {
	Id              int
	Fullname        string
	Email           string
	Username        sql.NullString
	Password        sql.NullString
	Role            string
	Image           []byte
	EmailVerifiedAt sql.NullTime
}
type CreateUserInput struct {
	Fullname string
	Image    []byte
}
