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

// MockCategoryRepository is a mock implementation of repository.CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) GetAll(ctx context.Context) ([]entities.CategoryEntity, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.CategoryEntity), args.Error(1)
}

func (m *MockCategoryRepository) GetById(ctx context.Context, id uint) (*entities.CategoryEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CategoryEntity), args.Error(1)
}

func (m *MockCategoryRepository) Create(ctx context.Context, request *entities.CategoryEntity) (*entities.CategoryEntity, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CategoryEntity), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, request *entities.CategoryEntity) (*entities.CategoryEntity, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CategoryEntity), args.Error(1)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCategoryService_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		expected := []entities.CategoryEntity{
			{ID: 1, Name: "Category 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		mockRepo.On("GetAll", mock.Anything).Return(expected, nil)

		res, appErr := svc.GetAll(context.Background())
		assert.Nil(t, appErr)
		assert.Len(t, res, 1)
		assert.Equal(t, "Category 1", res[0].Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		mockRepo.On("GetAll", mock.Anything).Return(nil, errors.New("db error"))

		res, appErr := svc.GetAll(context.Background())
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_GetById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		entity := &entities.CategoryEntity{ID: 1, Name: "Tech"}
		mockRepo.On("GetById", mock.Anything, uint(1)).Return(entity, nil)

		res, appErr := svc.GetById(context.Background(), 1)
		assert.Nil(t, appErr)
		assert.Equal(t, uint(1), res.ID)
		assert.Equal(t, "Tech", res.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		_, appErr := svc.GetById(context.Background(), 0)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		mockRepo.On("GetById", mock.Anything, uint(99)).Return(nil, repository.ErrCategoryNotFound)

		_, appErr := svc.GetById(context.Background(), 99)
		assert.NotNil(t, appErr)
		assert.Equal(t, 404, appErr.Code)
		assert.Equal(t, "category not found", appErr.Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		mockRepo.On("GetById", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

		_, appErr := svc.GetById(context.Background(), 1)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		req := payload.CategoryRequest{Name: "Fiction"}
		created := &entities.CategoryEntity{ID: 1, Name: "Fiction"}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *entities.CategoryEntity) bool {
			return c.Name == req.Name
		})).Return(created, nil)

		res, appErr := svc.Create(context.Background(), req)
		assert.Nil(t, appErr)
		assert.Equal(t, uint(1), res.ID)
		assert.Equal(t, "Fiction", res.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository failure", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		req := payload.CategoryRequest{Name: "Non-Fiction"}
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("insert failed"))

		_, appErr := svc.Create(context.Background(), req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		req := payload.CategoryRequest{Name: "Updated Tech"}
		existing := &entities.CategoryEntity{ID: 1, Name: "Old Tech"}
		updated := &entities.CategoryEntity{ID: 1, Name: "Updated Tech"}

		mockRepo.On("GetById", mock.Anything, uint(1)).Return(existing, nil)
		mockRepo.On("Update", mock.Anything, existing).Return(updated, nil)

		res, appErr := svc.Update(context.Background(), req, 1)
		assert.Nil(t, appErr)
		assert.Equal(t, "Updated Tech", res.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		_, appErr := svc.Update(context.Background(), payload.CategoryRequest{Name: "Test"}, 0)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		req := payload.CategoryRequest{Name: "Tech"}
		mockRepo.On("GetById", mock.Anything, uint(99)).Return(nil, repository.ErrCategoryNotFound)

		_, appErr := svc.Update(context.Background(), req, 99)
		assert.NotNil(t, appErr)
		assert.Equal(t, 404, appErr.Code)
		assert.Equal(t, "category not found", appErr.Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository update error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		req := payload.CategoryRequest{Name: "Tech"}
		existing := &entities.CategoryEntity{ID: 1, Name: "Tech"}

		mockRepo.On("GetById", mock.Anything, uint(1)).Return(existing, nil)
		mockRepo.On("Update", mock.Anything, existing).Return(nil, errors.New("db error"))

		_, appErr := svc.Update(context.Background(), req, 1)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

		appErr := svc.Delete(context.Background(), 1)
		assert.Nil(t, appErr)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		appErr := svc.Delete(context.Background(), 0)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		mockRepo.On("Delete", mock.Anything, uint(99)).Return(repository.ErrCategoryNotFound)

		appErr := svc.Delete(context.Background(), 99)
		assert.NotNil(t, appErr)
		assert.Equal(t, 404, appErr.Code)
		assert.Equal(t, "category not found", appErr.Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		svc := service.NewCategoryServiceImpl(mockRepo)

		mockRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("db error"))

		appErr := svc.Delete(context.Background(), 1)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}
