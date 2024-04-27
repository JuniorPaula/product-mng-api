package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"web_server/internal/dto"
	"web_server/internal/entity"
	"web_server/internal/infra/database"

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

func (h *UserHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpireIn := r.Context().Value("expires_in").(int)

	var user dto.JWTInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao decodificar o corpo da requisição"

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	u, err := h.UserDB.GetByEmail(user.Email)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Credenciais inválidas"

		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	if !u.ValidatePassword(user.Password) {
		paylod.Error = true
		paylod.Message = "Credenciais inválidas"

		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	_, token, _ := jwt.Encode(map[string]interface{}{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
		"exp":   time.Now().Add(time.Hour * time.Duration(jwtExpireIn)).Unix(),
	})

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: token,
	}

	paylod.Error = false
	paylod.Message = "sucesso"
	paylod.Data = accessToken

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
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

	u, err := entity.NewUser(input.Name, input.Email, input.Password)
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
