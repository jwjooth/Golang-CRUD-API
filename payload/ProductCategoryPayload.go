package payload

import (
	"time"

	"golang-restful-api/entities"
)

type ProductCategoryRequest struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required,min=2,max=2"`
}

type ProductCategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProductCategoryResponse(e entities.ProductCategoryEntity) ProductCategoryResponse {
	return ProductCategoryResponse{
		ID:        e.ID,
		Name:      e.Name,
		Code:      e.Code,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func NewProductCategoryResponses(list []entities.ProductCategoryEntity) []ProductCategoryResponse {
	out := make([]ProductCategoryResponse, 0, len(list))
	for _, e := range list {
		out = append(out, NewProductCategoryResponse(e))
	}
	return out
}
