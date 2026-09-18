package payload

import "time"

type BookRequest struct {
	Title      string `json:"title" validate:"required"`
	CategoryID uint   `json:"category_id" validate:"required"`
	Author     string `json:"author" validate:"required"`
	Stock      int    `json:"stock" validate:"required"`
}

type BookResponse struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	CategoryID uint      `json:"category_id"`
	Author     string    `json:"author"`
	Stock      int       `json:"stock"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ListBookMeta struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

type ListBookrResponse struct {
	Data []BookResponse `json:"data"`
	Meta ListBookMeta   `json:"meta"`
}
