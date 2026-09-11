package controller

import (
	"encoding/json"
	"golang-restful-api/helper"
	"golang-restful-api/payload"
	"golang-restful-api/service"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductControllerImpl struct {
	ProductService service.ProductService
}

func NewProductController(service service.ProductService) *ProductControllerImpl {
	return &ProductControllerImpl{
		ProductService: service,
	}
}

type ProductController interface {
	GetAllProduct(w http.ResponseWriter, r *http.Request)
	CreateProduct(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
}

func (p *ProductControllerImpl) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	result, appErr := p.ProductService.GetAllProduct(r.Context())
	if appErr != nil {
		helper.Logger(http.MethodGet, appErr.StatusCode, appErr.Message)
		w.WriteHeader(appErr.StatusCode)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": appErr.Message}); encErr != nil {
			log.Printf("failed to encode error response: %v", encErr)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (p *ProductControllerImpl) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request payload.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("POST /products bad request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encErr != nil {
			log.Printf("failed to encode error response: %v", encErr)
		}
		return
	}
	result, appErr := p.ProductService.CreateProduct(r.Context(), request)
	if appErr != nil {
		helper.Logger(http.MethodPost, appErr.StatusCode, appErr.Message)
		w.WriteHeader(appErr.StatusCode)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": appErr.Message}); encErr != nil {
			log.Printf("failed to encode error response: %v", encErr)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]any{"message": "success", "data": result}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (p *ProductControllerImpl) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("PUT /products/%s bad request: invalid id", idStr)
		w.WriteHeader(http.StatusBadRequest)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": "invalid product id"}); encErr != nil {
			log.Printf("failed to encode error response: %v", encErr)
		}
		return
	}

	var request payload.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("PUT /products/%d bad request: %v", id, err)
		w.WriteHeader(http.StatusBadRequest)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encErr != nil {
			log.Printf("failed to encode error response: %v", encErr)
		}
		return
	}

	result, appErr := p.ProductService.UpdateProduct(r.Context(), request, id)
	if appErr != nil {
		helper.Logger(http.MethodPut, appErr.StatusCode, appErr.Message)
		w.WriteHeader(appErr.StatusCode)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": appErr.Message}); encErr != nil {
			log.Printf("failed to encode error response: %v", encErr)
		}
		return
	}

	log.Printf("PUT /products/%d success", id)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]any{"message": "success", "data": result}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (p *ProductControllerImpl) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("DELETE /products/%s bad request: invalid id", idStr)
		w.WriteHeader(http.StatusBadRequest)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": "invalid product id"}); encErr != nil {
			log.Printf("failed to encode response: %v", encErr)
		}
		return
	}

	ok, appErr := p.ProductService.DeleteProduct(r.Context(), id)
	if appErr != nil {
		helper.Logger(http.MethodDelete, appErr.StatusCode, appErr.Message)
		w.WriteHeader(appErr.StatusCode)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": appErr.Message}); encErr != nil {
			log.Printf("failed to encode response: %v", encErr)
		}
		return
	}

	if !ok {
		log.Printf("DELETE /products/%d not found", id)
		w.WriteHeader(http.StatusNotFound)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": "product not found"}); encErr != nil {
			log.Printf("failed to encode response: %v", encErr)
		}
		return
	}

	log.Printf("DELETE /products/%d success", id)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "product deleted successfully"}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
