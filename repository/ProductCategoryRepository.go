package repository

import (
	"context"
	"errors"

	"golang-restful-api/entities"

	"gorm.io/gorm"
)

// ErrProductCategoryNotFound is returned when a product_category row does not exist.
var ErrProductCategoryNotFound = errors.New("product category not found")

// ProductCategoryRepository is a pure persistence boundary: entities in, entities out.
type ProductCategoryRepository interface {
	Create(ctx context.Context, request *entities.ProductCategoryEntity) (*entities.ProductCategoryEntity, error)
	GetAll(ctx context.Context) ([]entities.ProductCategoryEntity, error)
	GetById(ctx context.Context, id uint) (*entities.ProductCategoryEntity, error)
	Update(ctx context.Context, request *entities.ProductCategoryEntity) (*entities.ProductCategoryEntity, error)
	Delete(ctx context.Context, id uint) error
}

type productCategoryRepository struct {
	db *gorm.DB
}

// NewProductCategoryRepositoryImpl constructs a ProductCategoryRepository backed by GORM.
func NewProductCategoryRepositoryImpl(db *gorm.DB) ProductCategoryRepository {
	return &productCategoryRepository{db: db}
}

func (r *productCategoryRepository) Create(ctx context.Context, request *entities.ProductCategoryEntity) (*entities.ProductCategoryEntity, error) {
	if err := r.db.WithContext(ctx).Create(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (r *productCategoryRepository) GetAll(ctx context.Context) ([]entities.ProductCategoryEntity, error) {
	if err := r.db.WithContext(ctx).Model(&entities.ProductCategoryEntity{}).Error; err != nil {
		return nil, err
	}

	var records []entities.ProductCategoryEntity
	query := r.db.WithContext(ctx).Order("id asc")
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *productCategoryRepository) GetById(ctx context.Context, id uint) (*entities.ProductCategoryEntity, error) {
	var productCategory entities.ProductCategoryEntity
	if err := r.db.WithContext(ctx).First(&productCategory, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductCategoryNotFound
		}
		return nil, err
	}
	return &productCategory, nil
}

func (r *productCategoryRepository) Update(ctx context.Context, request *entities.ProductCategoryEntity) (*entities.ProductCategoryEntity, error) {
	if err := r.db.WithContext(ctx).Save(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (r *productCategoryRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entities.ProductCategoryEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrProductCategoryNotFound
	}
	return nil
}
