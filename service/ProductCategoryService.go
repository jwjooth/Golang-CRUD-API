package service

import (
	"context"
	"errors"

	"golang-restful-api/entities"
	"golang-restful-api/helper"
	"golang-restful-api/payload"
	"golang-restful-api/repository"

	"github.com/go-playground/validator/v10"
)

// ProductCategoryService holds business logic + validation and maps entities to DTOs.
type ProductCategoryService interface {
	GetAll(ctx context.Context) ([]payload.ProductCategoryResponse, *helper.AppError)
	GetById(ctx context.Context, id uint) (payload.ProductCategoryResponse, *helper.AppError)
	Create(ctx context.Context, req payload.ProductCategoryRequest) (payload.ProductCategoryResponse, *helper.AppError)
	Update(ctx context.Context, id uint, req payload.ProductCategoryRequest) (payload.ProductCategoryResponse, *helper.AppError)
	Delete(ctx context.Context, id uint) *helper.AppError
}

type productCategoryService struct {
	repository repository.ProductCategoryRepository
	validate   *validator.Validate
}

// NewProductCategoryServiceImpl constructs a ProductCategoryService.
func NewProductCategoryServiceImpl(repo repository.ProductCategoryRepository) ProductCategoryService {
	return &productCategoryService{repository: repo, validate: validator.New()}
}

// GetAll returns every product category as a response DTO.
func (s *productCategoryService) GetAll(ctx context.Context) ([]payload.ProductCategoryResponse, *helper.AppError) {
	records, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, helper.Internal("failed to fetch product categories", err)
	}
	return payload.NewProductCategoryResponses(records), nil
}

// GetById validates the ID and returns the matching product category.
func (s *productCategoryService) GetById(ctx context.Context, id uint) (payload.ProductCategoryResponse, *helper.AppError) {
	if id == 0 {
		return payload.ProductCategoryResponse{}, helper.BadRequest("invalid product category id")
	}
	productCategory, err := s.repository.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrProductCategoryNotFound) {
			return payload.ProductCategoryResponse{}, helper.NotFound("product category not found")
		}
		return payload.ProductCategoryResponse{}, helper.Internal("failed to fetch product category", err)
	}
	return payload.NewProductCategoryResponse(*productCategory), nil
}

// Create validates and persists a new product category.
func (s *productCategoryService) Create(ctx context.Context, req payload.ProductCategoryRequest) (payload.ProductCategoryResponse, *helper.AppError) {
	if err := s.validate.Struct(req); err != nil {
		return payload.ProductCategoryResponse{}, helper.BadRequest(helper.FormatValidationError(err))
	}

	productCategory := &entities.ProductCategoryEntity{
		Name: req.Name,
		Code: req.Code,
	}

	created, err := s.repository.Create(ctx, productCategory)
	if err != nil {
		return payload.ProductCategoryResponse{}, helper.Internal("failed to create product category", err)
	}
	return payload.NewProductCategoryResponse(*created), nil
}

// Update validates and persists changes to a product category.
func (s *productCategoryService) Update(ctx context.Context, id uint, req payload.ProductCategoryRequest) (payload.ProductCategoryResponse, *helper.AppError) {
	if id == 0 {
		return payload.ProductCategoryResponse{}, helper.BadRequest("invalid product category id")
	}
	if err := s.validate.Struct(req); err != nil {
		return payload.ProductCategoryResponse{}, helper.BadRequest(helper.FormatValidationError(err))
	}

	existing, err := s.repository.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrProductCategoryNotFound) {
			return payload.ProductCategoryResponse{}, helper.NotFound("product category not found")
		}
		return payload.ProductCategoryResponse{}, helper.Internal("failed to fetch product category", err)
	}

	existing.Name = req.Name
	existing.Code = req.Code

	updated, err := s.repository.Update(ctx, existing)
	if err != nil {
		return payload.ProductCategoryResponse{}, helper.Internal("failed to update product category", err)
	}
	return payload.NewProductCategoryResponse(*updated), nil
}

// Delete validates the ID and removes the matching product category.
func (s *productCategoryService) Delete(ctx context.Context, id uint) *helper.AppError {
	if id == 0 {
		return helper.BadRequest("invalid product category id")
	}
	if err := s.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrProductCategoryNotFound) {
			return helper.NotFound("product category not found")
		}
		return helper.Internal("failed to delete product category", err)
	}
	return nil
}
