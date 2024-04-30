package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"web_server/internal/dto"
	"web_server/internal/infra/database"

	"github.com/go-chi/jwtauth"
	"gorm.io/gorm"
)

type AuthHandler struct {
	UserDB database.UserInterface
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		UserDB: database.NewUser(db),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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
		"id":       u.ID,
		"email":    u.Email,
		"name":     u.Name,
		"is_admin": u.Admin,
		"exp":      time.Now().Add(time.Hour * time.Duration(jwtExpireIn)).Unix(),
	})

	data := struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		IsAdmin     bool   `json:"is_admin"`
		AccessToken string `json:"access_token"`
	}{
		ID:          u.ID.String(),
		Name:        u.Name,
		IsAdmin:     u.Admin,
		AccessToken: token,
	}

	paylod.Error = false
	paylod.Message = "Login realizado com sucesso"
	paylod.Data = data

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}

func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)

	var body struct {
		Token string `json:"token"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Erro ao decodificar o corpo da requisição"
		paylod.Data = nil

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	token, err := jwtauth.VerifyToken(jwt, body.Token)
	if err != nil {
		paylod.Error = true
		paylod.Message = "Token inválido"
		paylod.Data = nil

		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	if token == nil {
		paylod.Error = true
		paylod.Message = "Token inválido"
		paylod.Data = nil

		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	if token.Expiration().Before(time.Now()) {
		paylod.Error = true
		paylod.Message = "Token expirado"
		paylod.Data = nil

		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paylod)
		return
	}

	paylod.Error = false
	paylod.Message = "Token válido"
	paylod.Data = nil

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paylod)
}
