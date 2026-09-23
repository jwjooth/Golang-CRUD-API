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

type MockProductCategoryService struct {
	mock.Mock
}

// GetAll records a mocked request to list product categories.
func (m *MockProductCategoryService) GetAll(ctx context.Context) ([]payload.ProductCategoryResponse, *helper.AppError) {
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
	return args.Get(0).([]payload.ProductCategoryResponse), err
}

// GetById records a mocked request to retrieve a product category.
func (m *MockProductCategoryService) GetById(ctx context.Context, id uint) (payload.ProductCategoryResponse, *helper.AppError) {
	args := m.Called(ctx, id)
	var err *helper.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*helper.AppError)
	}
	return args.Get(0).(payload.ProductCategoryResponse), err
}

// Create records a mocked request to create a product category.
func (m *MockProductCategoryService) Create(ctx context.Context, request payload.ProductCategoryRequest) (payload.ProductCategoryResponse, *helper.AppError) {
	args := m.Called(ctx, request)
	var err *helper.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*helper.AppError)
	}
	return args.Get(0).(payload.ProductCategoryResponse), err
}

// Update records a mocked request to update a product category.
func (m *MockProductCategoryService) Update(ctx context.Context, id uint, request payload.ProductCategoryRequest) (payload.ProductCategoryResponse, *helper.AppError) {
	args := m.Called(ctx, id, request)
	var err *helper.AppError
	if args.Get(1) != nil {
		err = args.Get(1).(*helper.AppError)
	}
	return args.Get(0).(payload.ProductCategoryResponse), err
}

// Delete records a mocked request to delete a product category.
func (m *MockProductCategoryService) Delete(ctx context.Context, id uint) *helper.AppError {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*helper.AppError)
	}
	return nil
}

// setupProductCategoryRouter registers product category routes for controller tests.
func setupProductCategoryRouter(ctrl controller.ProductCategoryController) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/product-categories", ctrl.GetAll)
	r.Post("/product-categories", ctrl.Create)
	r.Get("/product-categories/{id}", ctrl.GetById)
	r.Put("/product-categories/{id}", ctrl.Update)
	r.Delete("/product-categories/{id}", ctrl.Delete)
	return r
}

func TestProductCategoryController_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		expected := []payload.ProductCategoryResponse{
			{ID: 1, Name: "Tech", Code: "T1"},
		}
		svc.On("GetAll", mock.Anything).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/product-categories", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Service error", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		svc.On("GetAll", mock.Anything).Return(nil, helper.Internal("db failure", nil))

		req := httptest.NewRequest(http.MethodGet, "/product-categories", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		svc.AssertExpectations(t)
	})
}

func TestProductCategoryController_GetById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		svc.On("GetById", mock.Anything, uint(1)).Return(payload.ProductCategoryResponse{ID: 1, Name: "Tech", Code: "T1"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/product-categories/1", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/product-categories/abc", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		svc.On("GetById", mock.Anything, uint(99)).Return(payload.ProductCategoryResponse{}, helper.NotFound("product category not found"))

		req := httptest.NewRequest(http.MethodGet, "/product-categories/99", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		svc.AssertExpectations(t)
	})
}

func TestProductCategoryController_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		createReq := payload.ProductCategoryRequest{Name: "Tech", Code: "T1"}
		svc.On("Create", mock.Anything, createReq).Return(payload.ProductCategoryResponse{ID: 2, Name: "Tech", Code: "T1"}, nil)

		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/product-categories", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid JSON body", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		req := httptest.NewRequest(http.MethodPost, "/product-categories", bytes.NewReader([]byte("{bad-json}")))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestProductCategoryController_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		updateReq := payload.ProductCategoryRequest{Name: "Updated Tech", Code: "T2"}
		svc.On("Update", mock.Anything, uint(1), updateReq).Return(payload.ProductCategoryResponse{ID: 1, Name: "Updated Tech", Code: "T2"}, nil)

		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, "/product-categories/1", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		req := httptest.NewRequest(http.MethodPut, "/product-categories/foo", bytes.NewReader([]byte("{}")))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		updateReq := payload.ProductCategoryRequest{Name: "Tech", Code: "T1"}
		svc.On("Update", mock.Anything, uint(404), updateReq).Return(payload.ProductCategoryResponse{}, helper.NotFound("product category not found"))

		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, "/product-categories/404", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		svc.AssertExpectations(t)
	})
}

func TestProductCategoryController_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		svc.On("Delete", mock.Anything, uint(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/product-categories/1", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		req := httptest.NewRequest(http.MethodDelete, "/product-categories/0", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		svc := new(MockProductCategoryService)
		ctrl := controller.NewProductCategoryController(svc)
		router := setupProductCategoryRouter(ctrl)

		svc.On("Delete", mock.Anything, uint(99)).Return(helper.NotFound("product category not found"))

		req := httptest.NewRequest(http.MethodDelete, "/product-categories/99", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		svc.AssertExpectations(t)
	})
}
