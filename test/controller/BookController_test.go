package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-restful-api/controller"
	"golang-restful-api/helper"
	"golang-restful-api/payload"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) GetAll(ctx context.Context, page, perPage int) ([]payload.BookResponse, int64, *helper.AppError) {
	args := m.Called(ctx, page, perPage)
	if args.Get(0) == nil {
		var err *helper.AppError
		if args.Get(2) != nil {
			err = args.Get(2).(*helper.AppError)
		}
		return nil, args.Get(1).(int64), err
	}
	var err *helper.AppError
	if args.Get(2) != nil {
		err = args.Get(2).(*helper.AppError)
	}
	return args.Get(0).([]payload.BookResponse), args.Get(1).(int64), err
}

func (m *MockBookService) GetById(ctx context.Context, id uint) (payload.BookResponse, *helper.AppError) {
	args := m.Called(ctx, id)
	var err *helper.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*helper.AppError)
	}
	return args.Get(0).(payload.BookResponse), err
}

func (m *MockBookService) Create(ctx context.Context, request payload.BookRequest) (payload.BookResponse, *helper.AppError) {
	args := m.Called(ctx, request)
	var err *helper.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*helper.AppError)
	}
	return args.Get(0).(payload.BookResponse), err
}

func (m *MockBookService) Update(ctx context.Context, request payload.BookRequest, id uint) (payload.BookResponse, *helper.AppError) {
	args := m.Called(ctx, request, id)
	var err *helper.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*helper.AppError)
	}
	return args.Get(0).(payload.BookResponse), err
}

func (m *MockBookService) Delete(ctx context.Context, id uint) *helper.AppError {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*helper.AppError)
	}
	return nil
}

func setupBookRouter(ctrl controller.BookController) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/books", ctrl.GetAll)
	r.Post("/books", ctrl.Create)
	r.Get("/books/{id}", ctrl.GetById)
	r.Put("/books/{id}", ctrl.Update)
	r.Delete("/books/{id}", ctrl.Delete)
	return r
}

func TestBookController_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		expected := []payload.BookResponse{
			{ID: 1, Title: "The Go Way", Author: "John", Stock: 5},
		}
		svc.On("GetAll", mock.Anything, 1, 10).Return(expected, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/books?page=1&per_page=10", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Service error", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		svc.On("GetAll", mock.Anything, 1, 10).Return(nil, int64(0), helper.Internal("db failure", nil))

		req := httptest.NewRequest(http.MethodGet, "/books?page=1&per_page=10", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		svc.AssertExpectations(t)
	})
}

func TestBookController_GetById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		svc.On("GetById", mock.Anything, uint(1)).Return(payload.BookResponse{ID: 1, Title: "Go Design Patterns"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/books/not-a-number", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		svc.On("GetById", mock.Anything, uint(404)).Return(payload.BookResponse{}, helper.NotFound("book not found"))

		req := httptest.NewRequest(http.MethodGet, "/books/404", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		svc.AssertExpectations(t)
	})
}

func TestBookController_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		createReq := payload.BookRequest{Title: "Domain-Driven Design", CategoryID: 1, Author: "Eric Evans", Stock: 10}
		svc.On("Create", mock.Anything, createReq).Return(payload.BookResponse{ID: 10, Title: "Domain-Driven Design"}, nil)

		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader([]byte("{invalid}")))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestBookController_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		updateReq := payload.BookRequest{Title: "Refactoring", CategoryID: 1, Author: "Martin Fowler", Stock: 7}
		svc.On("Update", mock.Anything, updateReq, uint(1)).Return(payload.BookResponse{ID: 1, Title: "Refactoring"}, nil)

		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		req := httptest.NewRequest(http.MethodPut, "/books/xyz", bytes.NewReader([]byte("{}")))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid JSON body", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		req := httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewReader([]byte("wrong")))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestBookController_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		svc.On("Delete", mock.Anything, uint(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/books/1", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		svc := new(MockBookService)
		ctrl := controller.NewBookControllerImpl(svc)
		router := setupBookRouter(ctrl)

		req := httptest.NewRequest(http.MethodDelete, "/books/abc", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
