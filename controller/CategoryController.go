package controller

import (
	"encoding/json"
	"golang-restful-api/helper"
	"golang-restful-api/payload"
	"golang-restful-api/service"
	"net/http"
)

type CategoryController interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type CategoryControllerImpl struct {
	Service service.CategoryService
}

func NewCategoryControllerImpl(svc service.CategoryService) *CategoryControllerImpl {
	return &CategoryControllerImpl{Service: svc}
}

// GetAll @Summary List all categories
// @Description Retrieve all categories
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {object} []payload.CategoryResponse
// @Failure 500 {object} helper.ErrorEnvelope
// @Router /categories [get]
func (c *CategoryControllerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	items, appErr := c.Service.GetAll(r.Context())
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "success", items)
}

// GetById @Summary Get a category by ID
// @Description Retrieve a single category by its unique ID
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} payload.CategoryResponse
// @Failure 400 {object} helper.ErrorEnvelope
// @Failure 404 {object} helper.ErrorEnvelope
// @Router /categories/{id} [get]
func (c *CategoryControllerImpl) GetById(w http.ResponseWriter, r *http.Request) {
	id, idErr := helper.ParseIDParam(r, "id")
	if idErr != nil {
		helper.WriteAppError(w, idErr)
		return
	}
	item, appErr := c.Service.GetById(r.Context(), id)
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "success", item)
}

// Create @Summary Create a new category
// @Description Create a new category with the provided details
// @Tags category
// @Accept json
// @Produce json
// @Param body body payload.CategoryRequest true "Category request payload"
// @Success 201 {object} payload.CategoryResponse
// @Failure 400 {object} helper.ErrorEnvelope
// @Failure 500 {object} helper.ErrorEnvelope
// @Router /categories [post]
func (c *CategoryControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	var request payload.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body:"+err.Error())
		return
	}
	result, appErr := c.Service.Create(r.Context(), request)
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	helper.WriteSuccess(w, http.StatusCreated, "success", result)
}

// Update @Summary Update an existing category
// @Description Update an existing category by its ID with the provided details
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param body body payload.CategoryRequest true "Category update payload"
// @Success 200 {object} payload.CategoryResponse
// @Failure 400 {object} helper.ErrorEnvelope
// @Failure 404 {object} helper.ErrorEnvelope
// @Failure 500 {object} helper.ErrorEnvelope
// @Router /categories/{id} [put]
func (c *CategoryControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	id, idErr := helper.ParseIDParam(r, "id")
	if idErr != nil {
		helper.WriteAppError(w, idErr)
		return
	}
	var request payload.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body:"+err.Error())
		return
	}
	res, err := c.Service.Update(r.Context(), request, id)
	if err != nil {
		helper.WriteAppError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "success", res)
}

// Delete @Summary Delete a category by ID
// @Description Delete a category by its unique ID
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} helper.ErrorEnvelope
// @Failure 404 {object} helper.ErrorEnvelope
// @Failure 500 {object} helper.ErrorEnvelope
// @Router /categories/{id} [delete]
func (c *CategoryControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	id, idErr := helper.ParseIDParam(r, "id")
	if idErr != nil {
		helper.WriteAppError(w, idErr)
		return
	}
	if err := c.Service.Delete(r.Context(), id); err != nil {
		helper.WriteAppError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "success", nil)

}
