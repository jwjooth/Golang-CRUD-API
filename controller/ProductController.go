package controller

import (
	"encoding/json"
	"net/http"

	"golang-restful-api/helper"
	"golang-restful-api/payload"
	"golang-restful-api/service"
)

// ProductController is the HTTP boundary: decode → service → encode.
// No business logic or SQL here.
type ProductController interface {
	ListProducts(w http.ResponseWriter, r *http.Request)
	GetProductByID(w http.ResponseWriter, r *http.Request)
	CreateProduct(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
}

type productController struct {
	service service.ProductService
}

// NewProductController constructs a ProductController.
func NewProductController(svc service.ProductService) ProductController {
	return &productController{service: svc}
}

// ListProducts handles GET /api/v1/products?page=&per_page=.
func (c *productController) ListProducts(w http.ResponseWriter, r *http.Request) {
	page := helper.QueryInt(r, "page", 1)
	perPage := helper.QueryInt(r, "per_page", 10)

	items, total, appErr := c.service.List(r.Context(), page, perPage)
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}

	totalPages := 0
	if perPage > 0 {
		totalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}
	meta := payload.ListProductsMeta{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}
	helper.WriteSuccessWithMeta(w, http.StatusOK, "success", items, meta)
}

// GetProductByID handles GET /api/v1/products/{id}.
func (c *productController) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, appErr := helper.ParseIDParam(r, "id")
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	result, err := c.service.GetByID(r.Context(), id)
	if err != nil {
		helper.WriteAppError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "success", result)
}

// CreateProduct handles POST /api/v1/products.
func (c *productController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req payload.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	result, appErr := c.service.Create(r.Context(), req)
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	helper.WriteSuccess(w, http.StatusCreated, "success", result)
}

// UpdateProduct handles PUT /api/v1/products/{id}.
func (c *productController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, appErr := helper.ParseIDParam(r, "id")
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	var req payload.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	result, err := c.service.Update(r.Context(), id, req)
	if err != nil {
		helper.WriteAppError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "success", result)
}

// DeleteProduct handles DELETE /api/v1/products/{id}.
func (c *productController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, appErr := helper.ParseIDParam(r, "id")
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	if err := c.service.Delete(r.Context(), id); err != nil {
		helper.WriteAppError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "product deleted successfully", nil)
}
