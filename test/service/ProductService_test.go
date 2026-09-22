package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang-restful-api/entities"
	"golang-restful-api/payload"
	"golang-restful-api/repository"
	"golang-restful-api/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of repository.ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProductEntity), args.Error(1)
}

func (m *MockProductRepository) FindAll(ctx context.Context, limit, offset int) ([]entities.ProductEntity, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]entities.ProductEntity), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) FindByID(ctx context.Context, id uint) (*entities.ProductEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProductEntity), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProductEntity), args.Error(1)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestProductService_List(t *testing.T) {
	t.Run("Success with default pagination", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		expectedList := []entities.ProductEntity{
			{ID: 1, Name: "Product 1", Price: 100, Stock: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, Name: "Product 2", Price: 200, Stock: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		mockRepo.On("FindAll", mock.Anything, 10, 0).Return(expectedList, int64(2), nil)

		res, total, appErr := svc.List(context.Background(), 1, 10)
		assert.Nil(t, appErr)
		assert.Equal(t, int64(2), total)
		assert.Len(t, res, 2)
		assert.Equal(t, "Product 1", res[0].Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Sanitize pagination params", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		mockRepo.On("FindAll", mock.Anything, 100, 0).Return([]entities.ProductEntity{}, int64(0), nil)

		_, _, appErr := svc.List(context.Background(), -1, 200)
		assert.Nil(t, appErr)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository returns error", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		mockRepo.On("FindAll", mock.Anything, 10, 0).Return(nil, int64(0), errors.New("db connection lost"))

		res, total, appErr := svc.List(context.Background(), 1, 10)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		entity := &entities.ProductEntity{ID: 1, Name: "Gadget", Price: 50, Stock: 2}
		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(entity, nil)

		res, appErr := svc.GetByID(context.Background(), 1)
		assert.Nil(t, appErr)
		assert.Equal(t, uint(1), res.ID)
		assert.Equal(t, "Gadget", res.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		_, appErr := svc.GetByID(context.Background(), 0)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
		assert.Equal(t, "invalid product id", appErr.Message)
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		mockRepo.On("FindByID", mock.Anything, uint(99)).Return(nil, repository.ErrNotFound)

		_, appErr := svc.GetByID(context.Background(), 99)
		assert.NotNil(t, appErr)
		assert.Equal(t, 404, appErr.Code)
		assert.Equal(t, "product not found", appErr.Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

		_, appErr := svc.GetByID(context.Background(), 1)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		req := payload.CreateProductRequest{
			Name:        "New Phone",
			Description: "Smartphone",
			Price:       999.99,
			Stock:       15,
		}

		createdEntity := &entities.ProductEntity{
			ID:          10,
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
		}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *entities.ProductEntity) bool {
			return p.Name == req.Name && p.Price == req.Price
		})).Return(createdEntity, nil)

		res, appErr := svc.Create(context.Background(), req)
		assert.Nil(t, appErr)
		assert.Equal(t, uint(10), res.ID)
		assert.Equal(t, "New Phone", res.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation failure - empty name", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		req := payload.CreateProductRequest{
			Name:  "",
			Price: 100,
			Stock: 5,
		}

		_, appErr := svc.Create(context.Background(), req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Validation failure - non-positive price", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		req := payload.CreateProductRequest{
			Name:  "Test",
			Price: 0,
			Stock: 5,
		}

		_, appErr := svc.Create(context.Background(), req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Validation failure - negative stock", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		req := payload.CreateProductRequest{
			Name:  "Test",
			Price: 100,
			Stock: -1,
		}

		_, appErr := svc.Create(context.Background(), req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Repository failure", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		req := payload.CreateProductRequest{
			Name:  "Test Item",
			Price: 10,
			Stock: 1,
		}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("insert failed"))

		_, appErr := svc.Create(context.Background(), req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		req := payload.UpdateProductRequest{
			Name:        "Updated Keyboard",
			Description: "Mechanical",
			Price:       120,
			Stock:       8,
		}

		existing := &entities.ProductEntity{ID: 1, Name: "Old Keyboard", Price: 100, Stock: 5}
		updated := &entities.ProductEntity{ID: 1, Name: req.Name, Description: req.Description, Price: req.Price, Stock: req.Stock}

		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(existing, nil)
		mockRepo.On("Update", mock.Anything, existing).Return(updated, nil)

		res, appErr := svc.Update(context.Background(), 1, req)
		assert.Nil(t, appErr)
		assert.Equal(t, "Updated Keyboard", res.Name)
		assert.Equal(t, 120.0, res.Price)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		_, appErr := svc.Update(context.Background(), 0, payload.UpdateProductRequest{Name: "Item", Price: 10, Stock: 1})
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Validation failure", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		_, appErr := svc.Update(context.Background(), 1, payload.UpdateProductRequest{Name: "", Price: 10})
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Product not found", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		req := payload.UpdateProductRequest{Name: "Item", Price: 10, Stock: 1}
		mockRepo.On("FindByID", mock.Anything, uint(2)).Return(nil, repository.ErrNotFound)

		_, appErr := svc.Update(context.Background(), 2, req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 404, appErr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update repo error", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		req := payload.UpdateProductRequest{Name: "Item", Price: 10, Stock: 1}
		existing := &entities.ProductEntity{ID: 1, Name: "Item", Price: 10, Stock: 1}
		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(existing, nil)
		mockRepo.On("Update", mock.Anything, existing).Return(nil, errors.New("db update error"))

		_, appErr := svc.Update(context.Background(), 1, req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

		appErr := svc.Delete(context.Background(), 1)
		assert.Nil(t, appErr)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		appErr := svc.Delete(context.Background(), 0)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		mockRepo.On("Delete", mock.Anything, uint(99)).Return(repository.ErrNotFound)

		appErr := svc.Delete(context.Background(), 99)
		assert.NotNil(t, appErr)
		assert.Equal(t, 404, appErr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete repo error", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		svc := service.NewProductServiceImpl(mockRepo)

		mockRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("db error"))

		appErr := svc.Delete(context.Background(), 1)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

// Test backward compatibility constructor
func TestNewProductServiceImplWithDB(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := service.NewProductServiceImplWithDB(nil, mockRepo)
	assert.NotNil(t, svc)
}
