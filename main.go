package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-restful-api/config"
	"golang-restful-api/controller"
	"golang-restful-api/helper"
	"golang-restful-api/repository"
	"golang-restful-api/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// .env is optional: missing file must not crash production.
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := config.OpenDB(cfg)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer config.CloseDB(db)

	if err := config.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		helper.WriteSuccess(w, http.StatusOK, "ok", map[string]string{"status": "up"})
	})

	//product
	productRepo := repository.NewProductRepositoryImpl(db)
	productService := service.NewProductServiceImpl(productRepo)
	productController := controller.NewProductController(productService)

	//book
	bookRepo := repository.NewBookRepositoryImpl(db)
	bookService := service.NewBookServiceImpl(bookRepo)
	bookController := controller.NewBookControllerImpl(bookService)

	//category
	categoryRepo := repository.NewCategoryRepositoryImpl(db)
	categoryService := service.NewCategoryServiceImpl(categoryRepo)
	categoryController := controller.NewCategoryControllerImpl(categoryService)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			r.Get("/", productController.ListProducts)
			r.Post("/", productController.CreateProduct)
			r.Get("/{id}", productController.GetProductByID)
			r.Put("/{id}", productController.UpdateProduct)
			r.Delete("/{id}", productController.DeleteProduct)
		})
		r.Route("/books", func(r chi.Router) {
			r.Get("/", bookController.GetAll)
			r.Post("/", bookController.Create)
			r.Get("/{id}", bookController.GetById)
			r.Put("/{id}", bookController.Update)
			r.Delete("/{id}", bookController.Delete)
		})
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", categoryController.GetAll)
			r.Post("/", categoryController.Create)
			r.Get("/{id}", categoryController.GetById)
			r.Put("/{id}", categoryController.Update)
			r.Delete("/{id}", categoryController.Delete)
		})
	})

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")
}
