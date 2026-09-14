package payload

import (
	"golang-restful-api/entities"
	"time"
)

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

func NewBookResponse(e entities.BookEntity) BookResponse {
	return BookResponse{
		ID:         e.ID,
		Title:      e.Title,
		CategoryID: *e.CategoryId,
		Author:     e.Author,
		Stock:      e.Stock,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

func NewBookResponses(list []entities.BookEntity) []BookResponse {
	out := make([]BookResponse, 0, len(list))
	for _, e := range list {
		out = append(out, NewBookResponse(e))
	}
	return out
}
