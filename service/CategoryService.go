package service

import (
	"errors"
	"golang-restful-api/payload"
	"context"
)

type CategoryService interface {
	GetAll(ctx context.Context) ([]payload.CategoryResponse, *helper.AppError)
	GetById(ctx context.Context, id uint) (payload.CategoryResponse, *helper.AppError)
	Create(ctx context.Context, request payload.CategoryRequest) (payload.CategoryResponse, *helper.AppError)
	Update(ctx context.Context, request payload.CategoryRequest, id uint) (payload.CategoryResponse, *helper.AppError)
	Delete(ctx context.Context, id uint) *helper.AppError
}

type CategoryServiceImpl struct {
	repository *repository.CategoryRepository
	validate *validator.Validate
}


