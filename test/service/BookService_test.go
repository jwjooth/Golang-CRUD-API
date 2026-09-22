package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang-restful-api/entities"
	"golang-restful-api/payload"
	"golang-restful-api/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBookRepository is a mock implementation of repository.BookRepository
type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) GetAll(ctx context.Context, limit, offset int) ([]entities.BookEntity, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]entities.BookEntity), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRepository) GetById(ctx context.Context, id uint) (*entities.BookEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.BookEntity), args.Error(1)
}

func (m *MockBookRepository) Create(ctx context.Context, request *entities.BookEntity) (*entities.BookEntity, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.BookEntity), args.Error(1)
}

func (m *MockBookRepository) Update(ctx context.Context, request *entities.BookEntity) (*entities.BookEntity, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.BookEntity), args.Error(1)
}

func (m *MockBookRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestBookService_GetAll(t *testing.T) {
	t.Run("Success with default pagination", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		catID := uint(1)
		expected := []entities.BookEntity{
			{ID: 1, Title: "Book A", CategoryId: &catID, Author: "Author A", Stock: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		mockRepo.On("GetAll", mock.Anything, 10, 0).Return(expected, int64(1), nil)

		res, total, appErr := svc.GetAll(context.Background(), 1, 10)
		assert.Nil(t, appErr)
		assert.Equal(t, int64(1), total)
		assert.Len(t, res, 1)
		assert.Equal(t, "Book A", res[0].Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Sanitize pagination params", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		mockRepo.On("GetAll", mock.Anything, 100, 0).Return([]entities.BookEntity{}, int64(0), nil)

		_, _, appErr := svc.GetAll(context.Background(), 0, 150)
		assert.Nil(t, appErr)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		mockRepo.On("GetAll", mock.Anything, 10, 0).Return(nil, int64(0), errors.New("query error"))

		res, total, appErr := svc.GetAll(context.Background(), 1, 10)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookService_GetById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		catID := uint(2)
		entity := &entities.BookEntity{ID: 1, Title: "Book A", CategoryId: &catID, Author: "Author A", Stock: 5}
		mockRepo.On("GetById", mock.Anything, uint(1)).Return(entity, nil)

		res, appErr := svc.GetById(context.Background(), 1)
		assert.Nil(t, appErr)
		assert.Equal(t, uint(1), res.ID)
		assert.Equal(t, "Book A", res.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		_, appErr := svc.GetById(context.Background(), 0)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Database error", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		mockRepo.On("GetById", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

		_, appErr := svc.GetById(context.Background(), 1)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		catID := uint(1)
		req := payload.BookRequest{
			Title:      "Go Programming",
			CategoryID: catID,
			Author:     "John Doe",
			Stock:      10,
		}

		created := &entities.BookEntity{
			ID:         5,
			Title:      req.Title,
			CategoryId: &catID,
			Author:     req.Author,
			Stock:      req.Stock,
		}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(b *entities.BookEntity) bool {
			return b.Title == req.Title && *b.CategoryId == req.CategoryID
		})).Return(created, nil)

		res, appErr := svc.Create(context.Background(), req)
		assert.Nil(t, appErr)
		assert.Equal(t, uint(5), res.ID)
		assert.Equal(t, "Go Programming", res.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation failure - missing fields", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		req := payload.BookRequest{
			Title: "",
		}

		_, appErr := svc.Create(context.Background(), req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Repository failure", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		req := payload.BookRequest{
			Title:      "Test Book",
			CategoryID: 1,
			Author:     "Author",
			Stock:      5,
		}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("insert failed"))

		_, appErr := svc.Create(context.Background(), req)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookService_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		catID := uint(1)
		req := payload.BookRequest{
			Title:      "Updated Title",
			CategoryID: catID,
			Author:     "Updated Author",
			Stock:      12,
		}

		existing := &entities.BookEntity{ID: 1, Title: "Old Title", CategoryId: &catID, Author: "Old Author", Stock: 5}
		updated := &entities.BookEntity{ID: 1, Title: req.Title, CategoryId: &catID, Author: req.Author, Stock: req.Stock}

		mockRepo.On("GetById", mock.Anything, uint(1)).Return(existing, nil)
		mockRepo.On("Update", mock.Anything, existing).Return(updated, nil)

		res, appErr := svc.Update(context.Background(), req, 1)
		assert.Nil(t, appErr)
		assert.Equal(t, "Updated Title", res.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		_, appErr := svc.Update(context.Background(), payload.BookRequest{Title: "T", CategoryID: 1, Author: "A", Stock: 1}, 0)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Validation failure", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		_, appErr := svc.Update(context.Background(), payload.BookRequest{Title: ""}, 1)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Repository update error", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		catID := uint(1)
		req := payload.BookRequest{Title: "Title", CategoryID: 1, Author: "Author", Stock: 1}
		existing := &entities.BookEntity{ID: 1, Title: "Title", CategoryId: &catID, Author: "Author", Stock: 1}

		mockRepo.On("GetById", mock.Anything, uint(1)).Return(existing, nil)
		mockRepo.On("Update", mock.Anything, existing).Return(nil, errors.New("update err"))

		_, appErr := svc.Update(context.Background(), req, 1)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookService_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

		appErr := svc.Delete(context.Background(), 1)
		assert.Nil(t, appErr)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID (0)", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		appErr := svc.Delete(context.Background(), 0)
		assert.NotNil(t, appErr)
		assert.Equal(t, 400, appErr.Code)
	})

	t.Run("Database error", func(t *testing.T) {
		mockRepo := new(MockBookRepository)
		svc := service.NewBookServiceImpl(mockRepo)

		mockRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("db delete error"))

		appErr := svc.Delete(context.Background(), 1)
		assert.NotNil(t, appErr)
		assert.Equal(t, 500, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}
