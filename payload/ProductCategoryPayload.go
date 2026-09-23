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

// NewProductCategoryResponse maps a product category entity to its response DTO.
func NewProductCategoryResponse(e entities.ProductCategoryEntity) ProductCategoryResponse {
	return ProductCategoryResponse{
		ID:        e.ID,
		Name:      e.Name,
		Code:      e.Code,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// NewProductCategoryResponses maps product category entities to response DTOs.
func NewProductCategoryResponses(list []entities.ProductCategoryEntity) []ProductCategoryResponse {
	out := make([]ProductCategoryResponse, 0, len(list))
	for _, e := range list {
		out = append(out, NewProductCategoryResponse(e))
	}
	return out
}
