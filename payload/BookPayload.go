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

type ListBookResponse struct {
	Data []BookResponse `json:"data"`
	Meta ListBookMeta   `json:"meta"`
}

// NewBookResponse maps a book entity to its response DTO.
func NewBookResponse(e entities.BookEntity) BookResponse {
	resp := BookResponse{
		ID:        e.ID,
		Title:     e.Title,
		Author:    e.Author,
		Stock:     e.Stock,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
	if e.CategoryId != nil {
		resp.CategoryID = *e.CategoryId
	}
	return resp
}

func NewBookResponses(list []entities.BookEntity) []BookResponse {
	out := make([]BookResponse, 0, len(list))
	for _, e := range list {
		out = append(out, NewBookResponse(e))
	}
	return out
}
