package repository

import (
	"context"
	"errors"
	"golang-restful-api/entities"

	"gorm.io/gorm"
)

type BookRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]entities.BookEntity, int64, error)
	Create(ctx context.Context, request *entities.BookEntity) (*entities.BookEntity, error)
	GetById(ctx context.Context, id uint) (*entities.BookEntity, error)
	Update(ctx context.Context, request *entities.BookEntity) (*entities.BookEntity, error)
	Delete(ctx context.Context, id uint) error
}

type BookRepositoryImpl struct {
	db *gorm.DB
}

func NewBookRepositoryImpl(db *gorm.DB) *BookRepositoryImpl {
	return &BookRepositoryImpl{db: db}
}

func (b *BookRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]entities.BookEntity, int64, error) {
	var total int64
	if err := b.db.WithContext(ctx).Model(&entities.BookEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []entities.BookEntity
	query := b.db.WithContext(ctx).Order("id asc")
	if limit > 0 {
		query.Limit(limit)
	}
	if offset > 0 {
		query.Offset(offset)
	}
	if err := query.Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

func (b *BookRepositoryImpl) GetById(ctx context.Context, id uint) (*entities.BookEntity, error) {
	var book entities.BookEntity
	if err := b.db.WithContext(ctx).First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (b *BookRepositoryImpl) Create(ctx context.Context, request *entities.BookEntity) (*entities.BookEntity, error) {
	if err := b.db.WithContext(ctx).Create(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (b *BookRepositoryImpl) Update(ctx context.Context, request *entities.BookEntity) (*entities.BookEntity, error) {
	if err := b.db.WithContext(ctx).Save(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (b *BookRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := b.db.WithContext(ctx).Delete(&entities.BookEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}
	return nil
}
