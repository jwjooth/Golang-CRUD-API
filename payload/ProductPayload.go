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
	Category    string    `json:"category"`
	ImageUrl    string    `json:"imageUrl"`
	Stock       int       `json:"stock"`
	Rating      float64   `json:"rating"`
	ReviewCount int       `json:"reviewCount"`
	SKU         string    `json:"sku"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductRequest is the inbound DTO for POST /products.
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description,omitempty" validate:"max=5000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	Category    string  `json:"category,omitempty" validate:"max=50"`
<<<<<<< HEAD
	ImageUrl    string  `json:"imageUrl,omitempty"`
	SKU         string  `json:"sku" validate:"required,max=50"`
||||||| parent of 5bab319 (fix: correct product persistence and validation)
	ImageUrl    string  `json:"imageUrl,omitempty"`
	SKU         string  `json:"sku,omitempty" validate:"max=50"`
=======
	ImageUrl    string  `json:"imageUrl,omitempty" validate:"omitempty,http_url,max=2048"`
	SKU         string  `json:"sku,omitempty" validate:"max=50"`
>>>>>>> 5bab319 (fix: correct product persistence and validation)
}

// UpdateProductRequest is the inbound DTO for PUT /products/{id} (full replace).
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description,omitempty" validate:"max=5000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	Category    string  `json:"category,omitempty" validate:"max=50"`
<<<<<<< HEAD
	ImageUrl    string  `json:"imageUrl,omitempty"`
	SKU         string  `json:"sku" validate:"required,max=50"`
||||||| parent of 5bab319 (fix: correct product persistence and validation)
	ImageUrl    string  `json:"imageUrl,omitempty"`
	SKU         string  `json:"sku,omitempty" validate:"max=50"`
=======
	ImageUrl    string  `json:"imageUrl,omitempty" validate:"omitempty,http_url,max=2048"`
	SKU         string  `json:"sku,omitempty" validate:"max=50"`
>>>>>>> 5bab319 (fix: correct product persistence and validation)
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
	sku := ""
	if e.SKU != nil {
		sku = *e.SKU
	}

	return ProductResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Price:       e.Price,
		Category:    e.Category,
		ImageUrl:    e.ImageUrl,
		Stock:       e.Stock,
		Rating:      e.Rating,
		ReviewCount: e.ReviewCount,
		SKU:         sku,
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
