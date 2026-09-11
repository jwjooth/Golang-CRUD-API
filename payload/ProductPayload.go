package payload

import (
	"time"

	"golang-restful-api/entities"
)

// ProductResponse is the outward DTO for a product.
type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductRequest is the inbound DTO for POST /products.
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description,omitempty" validate:"max=5000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
}

// UpdateProductRequest is the inbound DTO for PUT /products/{id} (full replace).
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description,omitempty" validate:"max=5000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
}

// ListProductsMeta carries pagination metadata.
type ListProductsMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ListProductsResponse is the paginated envelope for GET /products.
type ListProductsResponse struct {
	Data []ProductResponse `json:"data"`
	Meta ListProductsMeta  `json:"meta"`
}

// NewProductResponse maps an entity to its response DTO.
func NewProductResponse(e entities.ProductEntity) ProductResponse {
	return ProductResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Price:       e.Price,
		Stock:       e.Stock,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// NewProductResponses maps entities to response DTOs (never returns nil).
func NewProductResponses(list []entities.ProductEntity) []ProductResponse {
	out := make([]ProductResponse, 0, len(list))
	for _, e := range list {
		out = append(out, NewProductResponse(e))
	}
	return out
}
