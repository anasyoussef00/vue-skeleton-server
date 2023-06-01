package models

import (
	"database/sql"
	"time"
)

type Book struct {
	Id               int64          `json:"id" db:"id"`
	UserId           int64          `json:"userId" db:"user_id"`
	Title            string         `json:"title" db:"title"`
	AlternativeTitle sql.NullString `json:"alternativeTitle" db:"alternative_title"`
	Author           string         `json:"author" db:"author"`
	Description      string         `json:"description" db:"description"`
	ReleaseDate      time.Time      `json:"releaseDate" db:"release_date"`
	CreatedAt        time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time      `json:"updatedAt" db:"updated_at"`
	DeletedAt        sql.NullTime   `json:"deletedAt" db:"deleted_at"`
}

type BookReq struct {
	Title            string `json:"title"`
	AlternativeTitle string `json:"alternativeTitle"`
	Author           string `json:"author"`
	Description      string `json:"description"`
	ReleaseDate      string `json:"releaseDate"`
}

type BookResp struct {
	Id               int64    `json:"id" db:"id"`
	Title            string   `json:"title" db:"title"`
	AlternativeTitle *string  `json:"alternativeTitle" db:"alternative_title"`
	Author           string   `json:"author" db:"author"`
	Description      string   `json:"description" db:"description"`
	ReleaseDate      string   `json:"releaseDate" db:"release_date"`
	User             UserResp `json:"user" db:"user"`
}
