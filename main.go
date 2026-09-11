package main

import (
	"context"
	"golang-restful-api/config"
	"golang-restful-api/controller"
	"golang-restful-api/repository"
	"golang-restful-api/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := config.OpenDB()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	r := chi.NewRouter()

	productRepository := repository.NewProductRepositoryImpl(db)
	productService := service.NewProductServiceImpl(db, productRepository)
	productController := controller.NewProductController(productService)

	r.Route("/products", func(r chi.Router) {
		r.Get("/", productController.GetAllProduct)
		r.Post("/", productController.CreateProduct)
		r.Get("/{id}", productController.UpdateProduct)
		r.Get("/{id}", productController.DeleteProduct)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen: %s \n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")
}
