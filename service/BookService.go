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

type BookService interface {
	GetAll(ctx context.Context, page, perPage int) ([]payload.BookResponse, int64, *helper.AppError)
	GetById(ctx context.Context, id uint) (payload.BookResponse, *helper.AppError)
	Create(ctx context.Context, request payload.BookRequest) (payload.BookResponse, *helper.AppError)
	Update(ctx context.Context, request payload.BookRequest, id uint) (payload.BookResponse, *helper.AppError)
	Delete(ctx context.Context, id uint) *helper.AppError
}

type BookServiceImpl struct {
	repository repository.BookRepository
	validate   *validator.Validate
}

func NewBookServiceImpl(repo repository.BookRepository) BookService {
	return &BookServiceImpl{repository: repo, validate: validator.New()}
}

func (b *BookServiceImpl) GetAll(ctx context.Context, page, perPage int) ([]payload.BookResponse, int64, *helper.AppError) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	records, total, err := b.repository.GetAll(ctx, perPage, offset)
	if err != nil {
		return nil, 0, helper.Internal("failed to fetch books", err)
	}
	return payload.NewBookResponses(records), total, nil
}

func (b *BookServiceImpl) GetById(ctx context.Context, id uint) (payload.BookResponse, *helper.AppError) {
	if id == 0 {
		return payload.BookResponse{}, helper.BadRequest("invalid book id")
	}
	book, err := b.repository.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, errors.New("book not found")) {
			return payload.BookResponse{}, helper.NotFound("book not found")
		}
		return payload.BookResponse{}, helper.Internal("failed to fetch book", err)
	}
	return payload.NewBookResponse(*book), nil
}

func (b *BookServiceImpl) Create(ctx context.Context, request payload.BookRequest) (payload.BookResponse, *helper.AppError) {
	if err := b.validate.Struct(request); err != nil {
		return payload.BookResponse{}, helper.BadRequest(helper.FormatValidationError(err))
	}

	book := &entities.BookEntity{
		CategoryId: &request.CategoryID,
		Title:      request.Title,
		Author:     request.Author,
		Stock:      request.Stock,
	}

	created, err := b.repository.Create(ctx, book)
	if err != nil {
		return payload.BookResponse{}, helper.Internal("failed to create book", err)
	}
	return payload.NewBookResponse(*created), nil

}

func (b *BookServiceImpl) Update(ctx context.Context, request payload.BookRequest, id uint) (payload.BookResponse, *helper.AppError) {
	if id == 0 {
		return payload.BookResponse{}, helper.BadRequest("invalid book id")
	}
	if err := b.validate.Struct(request); err != nil {
		return payload.BookResponse{}, helper.BadRequest(helper.FormatValidationError(err))
	}
	existing, err := b.repository.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, errors.New("book not found")) {
			return payload.BookResponse{}, helper.NotFound("book not found")
		}
		return payload.BookResponse{}, helper.Internal("failed to fetch book", err)
	}
	existing.Title = request.Title
	existing.CategoryId = &request.CategoryID
	existing.Author = request.Author
	existing.Stock = request.Stock

	updated, err := b.repository.Update(ctx, existing)
	if err != nil {
		return payload.BookResponse{}, helper.Internal("failed to update book", err)
	}
	return payload.NewBookResponse(*updated), nil
}

func (b *BookServiceImpl) Delete(ctx context.Context, id uint) *helper.AppError {
	if id == 0 {
		return helper.BadRequest("invalid book id")
	}
	if err := b.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, errors.New("book not found")) {
			return helper.NotFound("book not found")
		}
		return helper.Internal("failed to delete book", err)
	}
	return nil
}
