package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"web_server/internal/infra/database/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProductHandler_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductDB := mock.NewMockProductInterface(ctrl)

	productHandler := NewProductHandler(mockProductDB)

	t.Run("should return error when decode body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/products", strings.NewReader("invalid body"))
		res := httptest.NewRecorder()

		productHandler.CreateProduct(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Erro ao decodificar o corpo da requisição")
	})
}
