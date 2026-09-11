package controller

import (
	"encoding/json"
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
	result, err := p.ProductService.GetAllProduct(r.Context())
	if err != nil {
		http.Error(w, err.Message, http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		if jsErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); jsErr != nil {
			log.Printf("failed to encode response: %v", jsErr)
		}
		return
	}
	result, err := p.ProductService.CreateProduct(r.Context(), request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if jsErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Message}); jsErr != nil {
			log.Printf("failed to encode response: %v", jsErr)
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
		w.WriteHeader(http.StatusBadRequest)
		if jsErr := json.NewEncoder(w).Encode(map[string]string{"error": "invalid product id"}); jsErr != nil {
			log.Printf("failed to encode response: %v", jsErr)
		}
		return
	}

	var request payload.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if jsErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); jsErr != nil {
			log.Printf("failed to encode response: %v", jsErr)
		}
		return
	}

	result, err := p.ProductService.UpdateProduct(r.Context(), request, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if jsErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); jsErr != nil {
			log.Printf("failed to encode response: %v", jsErr)
		}
		return
	}

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
		w.WriteHeader(http.StatusBadRequest)
		if jsErr := json.NewEncoder(w).Encode(map[string]string{"error": "invalid product id"}); jsErr != nil {
			log.Printf("failed to encode response: %v", jsErr)
		}
		return
	}

	_, err = p.ProductService.DeleteProduct(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if jsErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); jsErr != nil {
			log.Printf("failed to encode response: %v", jsErr)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "product deleted successfully"}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
