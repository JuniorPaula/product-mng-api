package main

import (
	"fmt"
	"net/http"
	"web_server/configs"
	"web_server/internal/entity"
	"web_server/internal/infra/database"
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
		productDB = database.NewProduct(db)

		productHandler = handlers.NewProductHandler(productDB)
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

	mux.Route("/users", func(r chi.Router) {
		r.Use(jwtauth.Verifier(cfg.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(isAdmin)

		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.GetUsers)
		r.Get("/{id}", userHandler.GetUser)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})

	fmt.Print("Server running on port :8000\n")
	http.ListenAndServe(":8000", mux)
}

// isAdmin is middleware that checks if the user is an admin
func isAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !claims["is_admin"].(bool) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
