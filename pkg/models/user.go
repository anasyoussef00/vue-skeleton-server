package models

import (
	"database/sql"
	"time"
)

type User struct {
	Id           int64        `json:"id" db:"id"`
	FirstName    string       `json:"firstName" db:"first_name"`
	LastName     string       `json:"lastName" db:"last_name"`
	Username     string       `json:"username" db:"username"`
	EmailAddress string       `json:"emailAddress" db:"email_address"`
	Password     string       `json:"password" db:"password"`
	Birthdate    time.Time    `json:"birthdate" db:"birthdate"`
	Gender       string       `json:"gender" db:"gender"`
	CreatedAt    time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time    `json:"updatedAt" db:"updated_at"`
	DeletedAt    sql.NullTime `json:"deletedAt" db:"deleted_at"`
}

type UserRegisterRequest struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Username     string `json:"username"`
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
	Birthdate    string `json:"birthdate"`
	Gender       string `json:"gender"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResp struct {
	Id           int64     `json:"id" db:"id"`
	FirstName    string    `json:"firstName" db:"first_name"`
	LastName     string    `json:"lastName" db:"last_name"`
	Username     string    `json:"username" db:"username"`
	EmailAddress string    `json:"emailAddress" db:"email_address"`
	Birthdate    time.Time `json:"birthdate" db:"birthdate"`
	Gender       string    `json:"gender" db:"gender"`
}

//func UserCtx(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		var user *User
//
//		if userID := chi.URLParam(r, "userID"); userID != "" {
//
//		}
//	})
//}
