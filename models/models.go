package models

import (
	"database/sql"
	"time"
)

type Models struct {
	DB DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Thread struct {
	ID             int            `json:"id"`
	Title          string         `json:"title"`
	Content        string         `json:"content"`
	AuthorID       int            `json:"author_id"`
	AuthorName     string         `json:"author_name"`
	Upvotes        int            `json:"upvotes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	IsSolved       bool           `json:"is_solved"`
	ThreadCategory map[int]string `json:"categories"`
	CategoryID     int            `json:"category_id"`
}

type Category struct {
	ID           int       `json:"id"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type ThreadCategory struct {
	ID         int       `json:"-"`
	ThreadID   int       `json:"-"`
	CategoryID int       `json:"-"`
	Category   Category  `json:"category"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Reply struct {
	ID         int       `json:"id"`
	ThreadID   int       `json:"thread_id"`
	Content    string    `json:"content"`
	AuthorID   int       `json:"author_id"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
	IsAnswer   bool      `json:"is_answer"`
}
