package handlers

import (
	"encoding/json"
	"net/http"
	"web_server/internal/dto"
	"web_server/internal/entity"
	"web_server/internal/infra/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
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

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		paylod.Error = true
		paylod.Message = "ID inválido"

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	user, err := h.UserDB.GetByID(id)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Usuário não encontrado"

		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	paylod.Error = false
	paylod.Message = "Usuário encontrado"
	paylod.Data = user

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	var input dto.UpdateUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao decodificar o corpo da requisição"
		paylod.Data = nil

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	u, err := h.UserDB.GetByID(id)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Usuário não encontrado"
		paylod.Data = nil

		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	// verify if the user wannas change the password
	if input.Password != "" {
		hash, err := u.GeneratePasswordHash(input.Password)
		if err != nil {
			paylod.Error = true
			paylod.Message = "Erro ao gerar hash da senha"
			paylod.Data = nil

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paylod)
			return
		}
		u.Password = hash
	}

	if input.Name != "" {
		u.Name = input.Name
	}
	if input.Email != "" {
		u.Email = input.Email
	}
	u.Admin = input.Admin

	err = h.UserDB.Update(u)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao atualizar usuário"
		paylod.Data = nil

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	paylod.Error = false
	paylod.Message = "Usuário atualizado"
	paylod.Data = u

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.UserDB.GetByID(id)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Usuário não encontrado"
		paylod.Data = nil

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())

	if claims["id"] == user.ID.String() {
		paylod.Error = true
		paylod.Message = "Você não pode deletar seu próprio usuário"
		paylod.Data = nil

		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	err = h.UserDB.Delete(id)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao deletar usuário"
		paylod.Data = nil

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	paylod.Error = false
	paylod.Message = "Usuário deletado"
	paylod.Data = nil

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}
