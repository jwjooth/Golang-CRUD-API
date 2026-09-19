package payload

import (
	"golang-restful-api/entities"
	"time"
)

type CategoryRequest struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCategoryResponse(e entities.CategoryEntity) CategoryResponse {
	return CategoryResponse{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func NewCategoryResponses(list []entities.CategoryEntity) []CategoryResponse {
	out := make([]CategoryResponse, 0, len(list))
	for _, e := range list {
		out = append(out, NewCategoryResponse(e))
	}
	return out
}
