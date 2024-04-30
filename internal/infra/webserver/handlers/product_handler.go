package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web_server/internal/dto"
	"web_server/internal/entity"
	"web_server/internal/infra/database"
	pkg "web_server/pkg/entity"

	"github.com/go-chi/chi/v5"
)

var paylod struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao decodificar o corpo da requisição"

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao criar produto"

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao criar produto"

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	paylod.Error = false
	paylod.Message = "Produto criado com sucesso"
	paylod.Data = p

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		paylod.Error = true
		paylod.Message = "ID inválido"

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	p, err := h.ProductDB.GetByID(id)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Produto não encontrado"

		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	paylod.Error = false
	paylod.Message = "Produto encontrado com sucesso"
	paylod.Data = p

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		paylod.Error = true
		paylod.Message = "ID inválido"
		paylod.Data = nil

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao decodificar o corpo da requisição"
		paylod.Data = nil

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	product.ID, err = pkg.ParseID(id)
	if err != nil {
		paylod.Error = true
		paylod.Message = "ID inválido"
		paylod.Data = nil

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	err = h.ProductDB.Update(&product)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao atualizar produto"
		paylod.Data = nil

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	paylod.Error = false
	paylod.Message = "Produto atualizado com sucesso"
	paylod.Data = product

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	paylod.Error = false
	paylod.Message = "Produto deletado com sucesso"

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	var (
		page  = r.URL.Query().Get("page")
		limit = r.URL.Query().Get("limit")
		sort  = r.URL.Query().Get("sort")
	)

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	products, err := h.ProductDB.GetAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	paylod.Error = false
	paylod.Message = "Produtos encontrados com sucesso"
	paylod.Data = products

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}
