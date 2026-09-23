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

// MockProductCategoryRepository is a mock implementation of repository.ProductCategoryRepository
type MockProductCategoryRepository struct {
	mock.Mock
}

// Create records a mocked request to persist a product category.
func (m *MockProductCategoryRepository) Create(ctx context.Context, request *entities.ProductCategoryEntity) (*entities.ProductCategoryEntity, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProductCategoryEntity), args.Error(1)
}

// GetAll records a mocked request to list product categories.
func (m *MockProductCategoryRepository) GetAll(ctx context.Context) ([]entities.ProductCategoryEntity, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.ProductCategoryEntity), args.Error(1)
}

// GetById records a mocked request to retrieve a product category.
func (m *MockProductCategoryRepository) GetById(ctx context.Context, id uint) (*entities.ProductCategoryEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProductCategoryEntity), args.Error(1)
}

// Update records a mocked request to persist product category changes.
func (m *MockProductCategoryRepository) Update(ctx context.Context, request *entities.ProductCategoryEntity) (*entities.ProductCategoryEntity, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProductCategoryEntity), args.Error(1)
}

// Delete records a mocked request to remove a product category.
func (m *MockProductCategoryRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestProductCategoryService_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		expected := []entities.ProductCategoryEntity{
			{ID: 1, Name: "Tech", Code: "T1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		mockRepo.On("GetAll", mock.Anything).Return(expected, nil)

		res, appErr := svc.GetAll(context.Background())
		assert.Nil(t, appErr)
		assert.Len(t, res, 1)
		assert.Equal(t, "Tech", res[0].Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		mockRepo.On("GetAll", mock.Anything).Return(nil, errors.New("db error"))

		res, appErr := svc.GetAll(context.Background())
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductCategoryService_GetById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		entity := &entities.ProductCategoryEntity{ID: 1, Name: "Tech", Code: "T1"}
		mockRepo.On("GetById", mock.Anything, uint(1)).Return(entity, nil)

		res, appErr := svc.GetById(context.Background(), 1)
		assert.Nil(t, appErr)
		assert.Equal(t, uint(1), res.ID)
		assert.Equal(t, "Tech", res.Name)
		assert.Equal(t, "T1", res.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		_, appErr := svc.GetById(context.Background(), 0)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		mockRepo.On("GetById", mock.Anything, uint(99)).Return(nil, repository.ErrProductCategoryNotFound)

		_, appErr := svc.GetById(context.Background(), 99)
		assert.NotNil(t, appErr)
		assert.Equal(t, 404, appErr.Code)
		assert.Equal(t, "product category not found", appErr.Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		mockRepo.On("GetById", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

		_, appErr := svc.GetById(context.Background(), 1)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductCategoryService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		req := payload.ProductCategoryRequest{Name: "Tech", Code: "T1"}
		created := &entities.ProductCategoryEntity{ID: 1, Name: req.Name, Code: req.Code}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *entities.ProductCategoryEntity) bool {
			return c.Name == req.Name && c.Code == req.Code
		})).Return(created, nil)

		res, appErr := svc.Create(context.Background(), req)
		assert.Nil(t, appErr)
		assert.Equal(t, uint(1), res.ID)
		assert.Equal(t, "Tech", res.Name)
		assert.Equal(t, "T1", res.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation failure - empty name", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		req := payload.ProductCategoryRequest{Name: "", Code: "T1"}

		_, appErr := svc.Create(context.Background(), req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Validation failure - code not 2 chars", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		req := payload.ProductCategoryRequest{Name: "Tech", Code: "T12"}

		_, appErr := svc.Create(context.Background(), req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Repository failure", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		req := payload.ProductCategoryRequest{Name: "Tech", Code: "T1"}
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("insert failed"))

		_, appErr := svc.Create(context.Background(), req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductCategoryService_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		req := payload.ProductCategoryRequest{Name: "Updated Tech", Code: "T2"}
		existing := &entities.ProductCategoryEntity{ID: 1, Name: "Old Tech", Code: "T1"}
		updated := &entities.ProductCategoryEntity{ID: 1, Name: req.Name, Code: req.Code}

		mockRepo.On("GetById", mock.Anything, uint(1)).Return(existing, nil)
		mockRepo.On("Update", mock.Anything, existing).Return(updated, nil)

		res, appErr := svc.Update(context.Background(), 1, req)
		assert.Nil(t, appErr)
		assert.Equal(t, "Updated Tech", res.Name)
		assert.Equal(t, "T2", res.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		_, appErr := svc.Update(context.Background(), 0, payload.ProductCategoryRequest{Name: "Test", Code: "T1"})
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		req := payload.ProductCategoryRequest{Name: "Tech", Code: "T1"}
		mockRepo.On("GetById", mock.Anything, uint(99)).Return(nil, repository.ErrProductCategoryNotFound)

		_, appErr := svc.Update(context.Background(), 99, req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 404, appErr.Code)
		assert.Equal(t, "product category not found", appErr.Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository update error", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		req := payload.ProductCategoryRequest{Name: "Tech", Code: "T1"}
		existing := &entities.ProductCategoryEntity{ID: 1, Name: "Tech", Code: "T1"}

		mockRepo.On("GetById", mock.Anything, uint(1)).Return(existing, nil)
		mockRepo.On("Update", mock.Anything, existing).Return(nil, errors.New("db error"))

		_, appErr := svc.Update(context.Background(), 1, req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductCategoryService_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

		appErr := svc.Delete(context.Background(), 1)
		assert.Nil(t, appErr)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		appErr := svc.Delete(context.Background(), 0)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		mockRepo.On("Delete", mock.Anything, uint(99)).Return(repository.ErrProductCategoryNotFound)

		appErr := svc.Delete(context.Background(), 99)
		assert.NotNil(t, appErr)
		assert.Equal(t, 404, appErr.Code)
		assert.Equal(t, "product category not found", appErr.Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		mockRepo := new(MockProductCategoryRepository)
		svc := service.NewProductCategoryServiceImpl(mockRepo)

		mockRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("db error"))

		appErr := svc.Delete(context.Background(), 1)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}
