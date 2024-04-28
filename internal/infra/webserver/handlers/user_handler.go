package handlers

import (
	"encoding/json"
	"net/http"
	"web_server/internal/dto"
	"web_server/internal/entity"
	"web_server/internal/infra/database"

	"gorm.io/gorm"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		UserDB: database.NewUser(db),
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao decodificar o corpo da requisição"

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	u, err := entity.NewUser(input.Name, input.Email, input.Password, input.Admin)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao criar usuário"

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	err = h.UserDB.Create(u)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao criar usuário"

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	paylod.Error = false
	paylod.Message = "Usuário criado com sucesso"

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserDB.GetAll()
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao buscar usuários"

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	paylod.Error = false
	paylod.Message = "Usuários encontrados"
	paylod.Data = users

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}
