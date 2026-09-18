package repository
<<<<<<< HEAD
||||||| eda2acc
=======

import (
	"context"
	"errors"
	"golang-restful-api/entities"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]entities.CategoryEntity, error)
	GetById(ctx context.Context, id uint) (*entities.CategoryEntity, error)
	Create(ctx context.Context, request *entities.CategoryEntity) (*entities.CategoryEntity, error)
	Update(ctx context.Context, request *entities.CategoryEntity) (*entities.CategoryEntity, error)
	Delete(ctx context.Context, id uint) error
}

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepositoryImpl(db *gorm.DB) *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{
		db: db,
	}
}

func (c *CategoryRepositoryImpl) GetAll(ctx context.Context) ([]entities.CategoryEntity, error) {
	if err := c.db.WithContext(ctx).Model(&entities.CategoryEntity{}).Error; err != nil {
		return nil, err
	}

	var records []entities.CategoryEntity
	query := c.db.WithContext(ctx).Order("id asc")
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (c *CategoryRepositoryImpl) GetById(ctx context.Context, id uint) (*entities.CategoryEntity, error) {
	var category entities.CategoryEntity
	if err := c.db.WithContext(ctx).First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
	}
	return &category, nil
}

func (c *CategoryRepositoryImpl) Create(ctx context.Context, request *entities.CategoryEntity) (*entities.CategoryEntity, error) {
	if err := c.db.WithContext(ctx).Create(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (c *CategoryRepositoryImpl) Update(ctx context.Context, request *entities.CategoryEntity) (*entities.CategoryEntity, error) {
	if err := c.db.WithContext(ctx).Save(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (c *CategoryRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := c.db.WithContext(ctx).Delete(&entities.CategoryEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}
>>>>>>> 95c71bc204f6f6f1eeb0c0cbc6697ec187e04c28
