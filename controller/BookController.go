package controller

import (
	"encoding/json"
	"golang-restful-api/helper"
	"golang-restful-api/payload"
	"golang-restful-api/service"
	"net/http"
)

type BookController interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type BookControllerImpl struct {
	service service.BookService
}

func NewBookControllerImpl(service service.BookService) BookController {
	return &BookControllerImpl{service: service}
}

func (b *BookControllerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	page := helper.QueryInt(r, "page", 1)
	perPage := helper.QueryInt(r, "per_page", 10)

	items, total, appErr := b.service.GetAll(r.Context(), page, perPage)
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}

	var totalPages int64
	if perPage > 0 {
		totalPages = (total + int64(perPage) - 1) / int64(perPage)
	}

	meta := payload.ListBookMeta{
		Page:      page,
		PerPage:   perPage,
		Total:     int(total),
		TotalPage: int(totalPages),
	}

	helper.WriteSuccessWithMeta(w, http.StatusOK, "success", items, meta)
}

func (b *BookControllerImpl) GetById(w http.ResponseWriter, r *http.Request) {
	id, appErr := helper.ParseIDParam(r, "id")
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	result, err := b.service.GetById(r.Context(), id)
	if err != nil {
		helper.WriteAppError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "success", result)
}

func (b *BookControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	var request payload.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	result, appErr := b.service.Create(r.Context(), request)
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	helper.WriteSuccess(w, http.StatusCreated, "success", result)
}

func (b *BookControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	id, appErr := helper.ParseIDParam(r, "id")
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	var request payload.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	result, err := b.service.Update(r.Context(), request, id)
	if err != nil {
		helper.WriteAppError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "success", result)
}

func (b *BookControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	id, appErr := helper.ParseIDParam(r, "id")
	if appErr != nil {
		helper.WriteAppError(w, appErr)
		return
	}
	if err := b.service.Delete(r.Context(), id); err != nil {
		helper.WriteAppError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "book deleted successfully", nil)
}
