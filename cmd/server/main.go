package main

import (
	"fmt"
	"net/http"
	"web_server/configs"
	"web_server/internal/entity"
	"web_server/internal/infra/webserver/handlers"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("teste.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	var (
		productHandler = handlers.NewProductHandler(db)
		userHandler    = handlers.NewUserHandler(db)
	)

	mux := chi.NewRouter()
	mux.Post("/products", productHandler.CreateProduct)
	mux.Get("/products", productHandler.GetProducts)
	mux.Get("/products/{id}", productHandler.GetProduct)
	mux.Put("/products/{id}", productHandler.UpdateProduct)
	mux.Delete("/products/{id}", productHandler.DeleteProduct)

	mux.Post("/users", userHandler.CreateUser)

	fmt.Print("Server running on port :8000\n")
	http.ListenAndServe(":8000", mux)
}
