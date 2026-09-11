package service

import (
	"context"
	"golang-restful-api/helper"
	"golang-restful-api/payload"
	"golang-restful-api/repository"

	"gorm.io/gorm"
)

type ProductService interface {
	GetAllProduct(ctx context.Context) ([]payload.ProductResponse, *helper.BaseErrorResponse)
	CreateProduct(ctx context.Context, request payload.ProductRequest) (payload.ProductResponse, *helper.BaseErrorResponse)
	UpdateProduct(ctx context.Context, request payload.ProductRequest, id int) (payload.ProductResponse, *helper.BaseErrorResponse)
	DeleteProduct(ctx context.Context, id int) (bool, *helper.BaseErrorResponse)
}

type ProductServiceImpl struct {
	db   *gorm.DB
	repo repository.ProductRepository
}

func NewProductServiceImpl(db *gorm.DB, repository repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		db:   db,
		repo: repository,
	}
}

func (p *ProductServiceImpl) GetAllProduct(ctx context.Context) ([]payload.ProductResponse, *helper.BaseErrorResponse) {
	result, err := p.repo.GetAllProduct(ctx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *ProductServiceImpl) CreateProduct(ctx context.Context, request payload.ProductRequest) (payload.ProductResponse, *helper.BaseErrorResponse) {
	result, err := p.repo.CreateProduct(ctx, request)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *ProductServiceImpl) UpdateProduct(ctx context.Context, request payload.ProductRequest, id int) (payload.ProductResponse, *helper.BaseErrorResponse) {
	result, err := p.repo.UpdateProduct(ctx, id, request)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *ProductServiceImpl) DeleteProduct(ctx context.Context, id int) (bool, *helper.BaseErrorResponse) {
	result, err := p.repo.DeleteProduct(ctx, id)
	if err != nil {
		return false, err
	}
	return result, nil
}
