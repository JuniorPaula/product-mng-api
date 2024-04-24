package handlers

import (
	"encoding/json"
	"net/http"
	"web_server/internal/dto"
	"web_server/internal/entity"
	"web_server/internal/infra/database"

	"gorm.io/gorm"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		ProductDB: database.NewProduct(db),
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
