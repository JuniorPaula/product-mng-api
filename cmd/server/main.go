package main

import (
	"fmt"
	"net/http"
	"web_server/configs"
	"web_server/internal/entity"
	"web_server/internal/infra/webserver/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg, err := configs.LoadConfig(".")
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
		authHandler    = handlers.NewAuthHandler(db)
	)

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.WithValue("jwt", cfg.TokenAuth))
	mux.Use(middleware.WithValue("expires_in", cfg.JWTExpiration))

	mux.Post("/login", authHandler.Login)

	mux.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(cfg.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	mux.Post("/users", userHandler.CreateUser)

	fmt.Print("Server running on port :8000\n")
	http.ListenAndServe(":8000", mux)
}
