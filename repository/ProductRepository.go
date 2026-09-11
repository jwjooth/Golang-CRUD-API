package repository

import (
	"context"
	"errors"

	"golang-restful-api/entities"

	"gorm.io/gorm"
)

// ErrNotFound is returned when a product row does not exist.
var ErrNotFound = errors.New("product not found")

// ProductRepository is a pure persistence boundary: entities in, entities out.
// No DTOs, no HTTP codes, no validation here.
type ProductRepository interface {
	Create(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error)
	FindAll(ctx context.Context, limit, offset int) ([]entities.ProductEntity, int64, error)
	FindByID(ctx context.Context, id uint) (*entities.ProductEntity, error)
	Update(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error)
	Delete(ctx context.Context, id uint) error
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepositoryImpl constructs a ProductRepository backed by GORM.
func NewProductRepositoryImpl(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error) {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) FindAll(ctx context.Context, limit, offset int) ([]entities.ProductEntity, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&entities.ProductEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []entities.ProductEntity
	q := r.db.WithContext(ctx).Order("id ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	if err := q.Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*entities.ProductEntity, error) {
	var product entities.ProductEntity
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error) {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&entities.ProductEntity{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
