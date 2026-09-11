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
		return payload.ProductResponse{}, errorHelper("name is required")
	}
	if request.Price < 0 {
		return payload.ProductResponse{}, errorHelper("price cant be less than zero")
	}
	if request.Price == 0 {
		return payload.ProductResponse{}, errorHelper("price is required")
	}
	if request.Stock < 0 {
		return payload.ProductResponse{}, errorHelper("stock cant be less than zero")
	}
	if request.Stock == 0 {
		return payload.ProductResponse{}, errorHelper("stock is required")
	}

	records := payload.ProductResponse{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	if err := p.db.WithContext(ctx).Create(&records).Error; err != nil {
		return payload.ProductResponse{}, errorHelper("failed to create product")
	}

	return records, nil
}

func (p *ProductRepositoryImpl) GetAllProduct(ctx context.Context) ([]payload.ProductResponse, *helper.BaseErrorResponse) {
	var records []payload.ProductResponse
	if dbErr := p.db.Find(&records).Error; dbErr != nil {
		return nil, errorHelper("failed to fetch products")
	}

	return records, nil
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
		return payload.ProductResponse{}, errorHelper("failed to update product data")
	}

	var updated entities.ProductEntity
	if err := p.db.WithContext(ctx).First(&updated, id).Error; err != nil {
		return payload.ProductResponse{}, errorHelper("failed to fetch updated product")
	}

	return payload.ProductResponse{
		Id:          updated.Id,
		Name:        updated.Name,
		Description: updated.Description,
		Price:       updated.Price,
		Stock:       updated.Stock,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}, nil
}

func (p *ProductRepositoryImpl) DeleteProduct(ctx context.Context, id int) (bool, *helper.BaseErrorResponse) {
	var records payload.ProductResponse

	err := p.db.Delete(&records).Error
	if err != nil {
		return false, errorHelper("failed to delete product data")
	}

	return true, nil
}

func errorHelper(kind string) *helper.BaseErrorResponse {
	return &helper.BaseErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    kind,
		Data:       nil,
	}
}
