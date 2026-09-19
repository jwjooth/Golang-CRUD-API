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

type CategoryService interface {
	GetAll(ctx context.Context) ([]payload.CategoryResponse, *helper.AppError)
	GetById(ctx context.Context, id uint) (payload.CategoryResponse, *helper.AppError)
	Create(ctx context.Context, request payload.CategoryRequest) (payload.CategoryResponse, *helper.AppError)
	Update(ctx context.Context, request payload.CategoryRequest, id uint) (payload.CategoryResponse, *helper.AppError)
	Delete(ctx context.Context, id uint) *helper.AppError
}

type CategoryServiceImpl struct {
	repository repository.CategoryRepository
	validate   *validator.Validate
}

func NewCategoryServiceImpl(repo repository.CategoryRepository) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		repository: repo,
		validate:   validator.New(),
	}
}

func (c *CategoryServiceImpl) GetAll(ctx context.Context) ([]payload.CategoryResponse, *helper.AppError) {
	records, err := c.repository.GetAll(ctx)
	if err != nil {
		return nil, helper.Internal("failed to fetch categories", err)
	}
	return payload.NewCategoryResponses(records), nil
}

func (c *CategoryServiceImpl) GetById(ctx context.Context, id uint) (payload.CategoryResponse, *helper.AppError) {
	if id <= 0 {
		return payload.CategoryResponse{}, helper.BadRequest("invalid category id")
	}
	category, err := c.repository.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, errors.New("category not found")) {
			return payload.CategoryResponse{}, helper.NotFound("category not found")
		}
		return payload.CategoryResponse{}, helper.Internal("failed to fetch category", err)
	}
	return payload.NewCategoryResponse(*category), nil
}

func (c *CategoryServiceImpl) Create(ctx context.Context, request payload.CategoryRequest) (payload.CategoryResponse, *helper.AppError) {
	if err := c.validate.Struct(request); err != nil {
		return payload.CategoryResponse{}, helper.BadRequest(helper.FormatValidationError(err))
	}

	ctgr := &entities.CategoryEntity{
		Name: request.Name,
	}

	category, err := c.repository.Create(ctx, ctgr)
	if err != nil {
		return payload.CategoryResponse{}, helper.Internal("failed to create category", err)
	}

	return payload.NewCategoryResponse(*category), nil
}

func (c *CategoryServiceImpl) Update(ctx context.Context, request payload.CategoryRequest, id uint) (payload.CategoryResponse, *helper.AppError) {
	if id <= 0 {
		return payload.CategoryResponse{}, helper.BadRequest("invalid category id")
	}
	if err := c.validate.Struct(request); err != nil {
		return payload.CategoryResponse{}, helper.BadRequest(helper.FormatValidationError(err))
	}

	existing, err := c.repository.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, errors.New("category not found")) {
			return payload.CategoryResponse{}, helper.NotFound("category not found")
		}
		return payload.CategoryResponse{}, helper.Internal("failed to fetch category", err)
	}
	existing.Name = request.Name

	category, err := c.repository.Update(ctx, existing)
	if err != nil {
		return payload.CategoryResponse{}, helper.Internal("failed to update category", err)
	}

	return payload.NewCategoryResponse(*category), nil
}

func (c *CategoryServiceImpl) Delete(ctx context.Context, id uint) *helper.AppError {
	if id <= 0 {
		return helper.BadRequest("invalid category id")
	}
	if err := c.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, errors.New("category not found")) {
			return helper.NotFound("category not found")
		}
		return helper.Internal("failed to delete category", err)
	}
	return nil
}
