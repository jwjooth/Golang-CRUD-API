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

type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) GetAll(ctx context.Context) ([]payload.CategoryResponse, *helper.AppError) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		var err *helper.AppError
		if args.Get(1) != nil {
			err = args.Get(1).(*helper.AppError)
		}
		return nil, err
	}
	var err *helper.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*helper.AppError)
	}
	return args.Get(0).([]payload.CategoryResponse), err
}

func (m *MockCategoryService) GetById(ctx context.Context, id uint) (payload.CategoryResponse, *helper.AppError) {
	args := m.Called(ctx, id)
	var err *helper.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*helper.AppError)
	}
	return args.Get(0).(payload.CategoryResponse), err
}

func (m *MockCategoryService) Create(ctx context.Context, request payload.CategoryRequest) (payload.CategoryResponse, *helper.AppError) {
	args := m.Called(ctx, request)
	var err *helper.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*helper.AppError)
	}
	return args.Get(0).(payload.CategoryResponse), err
}

func (m *MockCategoryService) Update(ctx context.Context, request payload.CategoryRequest, id uint) (payload.CategoryResponse, *helper.AppError) {
	args := m.Called(ctx, request, id)
	var err *helper.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*helper.AppError)
	}
	return args.Get(0).(payload.CategoryResponse), err
}

func (m *MockCategoryService) Delete(ctx context.Context, id uint) *helper.AppError {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*helper.AppError)
	}
	return nil
}

func setupCategoryRouter(ctrl controller.CategoryController) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/categories", ctrl.GetAll)
	r.Post("/categories", ctrl.Create)
	r.Get("/categories/{id}", ctrl.GetById)
	r.Put("/categories/{id}", ctrl.Update)
	r.Delete("/categories/{id}", ctrl.Delete)
	return r
}

func TestCategoryController_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockCategoryService)
		ctrl := controller.NewCategoryControllerImpl(svc)
		router := setupCategoryRouter(ctrl)

		expected := []payload.CategoryResponse{
			{ID: 1, Name: "Fiction"},
		}
		svc.On("GetAll", mock.Anything).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Service error", func(t *testing.T) {
		svc := new(MockCategoryService)
		ctrl := controller.NewCategoryControllerImpl(svc)
		router := setupCategoryRouter(ctrl)

		svc.On("GetAll", mock.Anything).Return(nil, helper.Internal("db failure", nil))

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		svc.AssertExpectations(t)
	})
}

func TestCategoryController_GetById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockCategoryService)
		ctrl := controller.NewCategoryControllerImpl(svc)
		router := setupCategoryRouter(ctrl)

		svc.On("GetById", mock.Anything, uint(1)).Return(payload.CategoryResponse{ID: 1, Name: "Comics"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		svc := new(MockCategoryService)
		ctrl := controller.NewCategoryControllerImpl(svc)
		router := setupCategoryRouter(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/categories/abc", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestCategoryController_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockCategoryService)
		ctrl := controller.NewCategoryControllerImpl(svc)
		router := setupCategoryRouter(ctrl)

		createReq := payload.CategoryRequest{Name: "Drama"}
		svc.On("Create", mock.Anything, createReq).Return(payload.CategoryResponse{ID: 5, Name: "Drama"}, nil)

		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid JSON body", func(t *testing.T) {
		svc := new(MockCategoryService)
		ctrl := controller.NewCategoryControllerImpl(svc)
		router := setupCategoryRouter(ctrl)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader([]byte("{bad-json}")))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestCategoryController_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockCategoryService)
		ctrl := controller.NewCategoryControllerImpl(svc)
		router := setupCategoryRouter(ctrl)

		updateReq := payload.CategoryRequest{Name: "Romance"}
		svc.On("Update", mock.Anything, updateReq, uint(1)).Return(payload.CategoryResponse{ID: 1, Name: "Romance"}, nil)

		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		svc := new(MockCategoryService)
		ctrl := controller.NewCategoryControllerImpl(svc)
		router := setupCategoryRouter(ctrl)

		req := httptest.NewRequest(http.MethodPut, "/categories/foo", bytes.NewReader([]byte("{}")))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestCategoryController_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockCategoryService)
		ctrl := controller.NewCategoryControllerImpl(svc)
		router := setupCategoryRouter(ctrl)

		svc.On("Delete", mock.Anything, uint(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		svc := new(MockCategoryService)
		ctrl := controller.NewCategoryControllerImpl(svc)
		router := setupCategoryRouter(ctrl)

		req := httptest.NewRequest(http.MethodDelete, "/categories/0", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
