package repository

import (
	"context"
	"golang-restful-api/entities"
	"golang-restful-api/helper"
	"golang-restful-api/payload"
	"net/http"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, request payload.ProductRequest) (payload.ProductResponse, *helper.BaseErrorResponse)
	GetAllProduct(ctx context.Context) ([]payload.ProductResponse, *helper.BaseErrorResponse)
	UpdateProduct(ctx context.Context, id int, request payload.ProductRequest) (payload.ProductResponse, *helper.BaseErrorResponse)
	DeleteProduct(ctx context.Context, id int) (bool, *helper.BaseErrorResponse)
}

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepositoryImpl(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{
		db: db,
	}
}

func (p *ProductRepositoryImpl) CreateProduct(ctx context.Context, request payload.ProductRequest) (payload.ProductResponse, *helper.BaseErrorResponse) {
	if request.Name == "" {
		return payload.ProductResponse{}, errorHelper(http.StatusBadRequest, "name is required")
	}
	if request.Price < 0 {
		return payload.ProductResponse{}, errorHelper(http.StatusBadRequest, "price cant be less than zero")
	}
	if request.Price == 0 {
		return payload.ProductResponse{}, errorHelper(http.StatusBadRequest, "price is required")
	}
	if request.Stock < 0 {
		return payload.ProductResponse{}, errorHelper(http.StatusBadRequest, "stock cant be less than zero")
	}
	if request.Stock == 0 {
		return payload.ProductResponse{}, errorHelper(http.StatusBadRequest, "stock is required")
	}

	product := entities.ProductEntity{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	if err := p.db.WithContext(ctx).Create(&product).Error; err != nil {
		return payload.ProductResponse{}, errorHelper(http.StatusInternalServerError, "failed to create product")
	}

	return toResponse(product), nil
}

func (p *ProductRepositoryImpl) GetAllProduct(ctx context.Context) ([]payload.ProductResponse, *helper.BaseErrorResponse) {
	var records []entities.ProductEntity
	if dbErr := p.db.WithContext(ctx).Find(&records).Error; dbErr != nil {
		return nil, errorHelper(http.StatusInternalServerError, "failed to fetch products")
	}

	return toResponses(records), nil
}

func (p *ProductRepositoryImpl) UpdateProduct(ctx context.Context, id int, request payload.ProductRequest) (payload.ProductResponse, *helper.BaseErrorResponse) {
	updates := make(map[string]any)

	if request.Name != "" {
		updates["name"] = request.Name
	}
	if request.Price != 0 {
		updates["price"] = request.Price
	}
	if request.Stock != 0 {
		updates["stock"] = request.Stock
	}

	if dbErr := p.db.WithContext(ctx).Model(&entities.ProductEntity{}).Where("id = ?", id).Updates(updates).Error; dbErr != nil {
		return payload.ProductResponse{}, errorHelper(http.StatusInternalServerError, "failed to update product data")
	}

	var updated entities.ProductEntity
	if err := p.db.WithContext(ctx).First(&updated, id).Error; err != nil {
		return payload.ProductResponse{}, errorHelper(http.StatusInternalServerError, "failed to fetch updated product")
	}

	return toResponse(updated), nil
}

func (p *ProductRepositoryImpl) DeleteProduct(ctx context.Context, id int) (bool, *helper.BaseErrorResponse) {
	if err := p.db.WithContext(ctx).Delete(&entities.ProductEntity{}, id).Error; err != nil {
		return false, errorHelper(http.StatusInternalServerError, "failed to delete product data")
	}

	return true, nil
}

func errorHelper(statusCode int, message string) *helper.BaseErrorResponse {
	return &helper.BaseErrorResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
	}
}

func toResponse(product entities.ProductEntity) payload.ProductResponse {
	return payload.ProductResponse{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func toResponses(products []entities.ProductEntity) []payload.ProductResponse {
	resp := make([]payload.ProductResponse, len(products))
	for i, p := range products {
		resp[i] = toResponse(p)
	}
	return resp
}
