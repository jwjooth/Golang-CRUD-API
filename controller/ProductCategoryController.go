package controller

import (
	"encoding/json"
	"net/http"

	"golang-restful-api/helper"
	"golang-restful-api/payload"
	"golang-restful-api/service"
)

// ProductCategoryController is the HTTP boundary: decode → service → encode.
// No business logic or SQL here.
type ProductCategoryController interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type productCategoryController struct {
	service service.ProductCategoryService
}

// NewProductCategoryController constructs a ProductCategoryController.
func NewProductCategoryController(svc service.ProductCategoryService) ProductCategoryController {
	return &productCategoryController{service: svc}
}

// GetAll @Summary List all product categories
// @Description Retrieve a list of all product categories
// @Tags product-category
// @Accept json
// @Produce json
// @Success 200 {object} []payload.ProductCategoryResponse
// @Failure 500 {object} helper.ErrorEnvelope
// @Router /product-categories [get]
func (c *productCategoryController) GetAll(w http.ResponseWriter, r *http.Request) {
	items, appErr := c.service.GetAll(r.Context())
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "success", items)
}

// GetById @Summary Get a product category by ID
// @Description Retrieve a single product category by its unique ID
// @Tags product-category
// @Accept json
// @Produce json
// @Param id path int true "Product Category ID"
// @Success 200 {object} payload.ProductCategoryResponse
// @Failure 400 {object} helper.ErrorEnvelope
// @Failure 404 {object} helper.ErrorEnvelope
// @Router /product-categories/{id} [get]
func (c *productCategoryController) GetById(w http.ResponseWriter, r *http.Request) {
	id, appErr := helper.ParseIDParam(r, "id")
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	result, err := c.service.GetById(r.Context(), id)
	if err != nil {
		helper.WriteAppError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "success", result)
}

// Create @Summary Create a new product category
// @Description Create a new product category with the provided details
// @Tags product-category
// @Accept json
// @Produce json
// @Param body body payload.ProductCategoryRequest true "Product category request payload"
// @Success 201 {object} payload.ProductCategoryResponse
// @Failure 400 {object} helper.ErrorEnvelope
// @Failure 500 {object} helper.ErrorEnvelope
// @Router /product-categories [post]
func (c *productCategoryController) Create(w http.ResponseWriter, r *http.Request) {
	var req payload.ProductCategoryRequest
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

// Update @Summary Update an existing product category
// @Description Update an existing product category by its ID with the provided details
// @Tags product-category
// @Accept json
// @Produce json
// @Param id path int true "Product Category ID"
// @Param body body payload.ProductCategoryRequest true "Product category update payload"
// @Success 200 {object} payload.ProductCategoryResponse
// @Failure 400 {object} helper.ErrorEnvelope
// @Failure 404 {object} helper.ErrorEnvelope
// @Failure 500 {object} helper.ErrorEnvelope
// @Router /product-categories/{id} [put]
func (c *productCategoryController) Update(w http.ResponseWriter, r *http.Request) {
	id, appErr := helper.ParseIDParam(r, "id")
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	var req payload.ProductCategoryRequest
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

// Delete @Summary Delete a product category by ID
// @Description Delete a product category by its unique ID
// @Tags product-category
// @Accept json
// @Produce json
// @Param id path int true "Product Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} helper.ErrorEnvelope
// @Failure 404 {object} helper.ErrorEnvelope
// @Failure 500 {object} helper.ErrorEnvelope
// @Router /product-categories/{id} [delete]
func (c *productCategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	id, appErr := helper.ParseIDParam(r, "id")
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	if err := c.service.Delete(r.Context(), id); err != nil {
		helper.WriteAppError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "product category deleted successfully", nil)
}
