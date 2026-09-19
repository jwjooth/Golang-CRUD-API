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

func (c *CategoryControllerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	items, appErr := c.Service.GetAll(r.Context())
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "success", items)
}

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
