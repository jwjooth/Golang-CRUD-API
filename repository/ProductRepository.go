package repository

import (
	"golang-restful-api/helper"
	"golang-restful-api/payload"
	"net/http"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(db *gorm.DB, request payload.ProductRequest) (payload.ProductResponse, *helper.BaseErrorResponse)
	GetAllProduct(db *gorm.DB) ([]payload.ProductResponse, *helper.BaseErrorResponse)
	UpdateProduct(db *gorm.DB, id int, request payload.ProductRequest) (payload.ProductResponse, *helper.BaseErrorResponse)
	DeleteProduct(db *gorm.DB, id int) bool
}

type ProductRepositoryImpl struct {
}

func (p *ProductRepositoryImpl) CreateProduct(db *gorm.DB, request payload.ProductRequest) (payload.ProductResponse, *helper.BaseErrorResponse) {
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

	return records, nil
}

func (p *ProductRepositoryImpl) GetAllProduct(db *gorm.DB) ([]payload.ProductResponse, *helper.BaseErrorResponse) {
	var records []payload.ProductResponse
	query := db.Select("select * from products")

	if dbErr := query.Scan(&records).Error; dbErr != nil {
		return nil, &helper.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get all product data",
			Data:       nil,
		}
	}

	return records, nil
}

func (p *ProductRepositoryImpl) UpdateProduct(db *gorm.DB, id int, request payload.ProductRequest) (payload.ProductResponse, *helper.BaseErrorResponse) {
	var records payload.ProductResponse

	if request.Name != "" {
		records.Name = request.Name
	}
	if request.Price != 0 {
		records.Price = request.Price
	}
	if request.Stock != 0 {
		records.Stock = request.Stock
	}

	if dbErr := db.Where("id = ?", id).First(&records).Error; dbErr != nil {
		return payload.ProductResponse{}, errorHelper("failed to update product data")
	}

	return records, nil
}

func (p *ProductRepositoryImpl) DeleteProduct(db *gorm.DB, id int) bool {
	var records payload.ProductResponse

	err := db.Delete(&records).Error
	if err != nil {
		return false
	}

	return true
}

func errorHelper(kind string) *helper.BaseErrorResponse {
	return &helper.BaseErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    kind,
		Data:       nil,
	}
}
