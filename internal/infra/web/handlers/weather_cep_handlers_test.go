package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

type Temperatura struct {
	TempC float64 `json:"temp_C"` // Temperatura em Celsius
	TempF float64 `json:"temp_F"` // Temperatura em Fahrenheit
	TempK float64 `json:"temp_K"` // Temperatura em Kelvin
}

// HTTP CODE 200  - Return body
func TestHandlerClimaCode200(t *testing.T) {
	router := chi.NewRouter()
	router.Get("/{cep}", HandlerClima)
	req, err := http.NewRequest("GET", "/34012690", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NoError(t, err)

	var Response Temperatura
	err = json.Unmarshal(w.Body.Bytes(), &Response)
	assert.NoError(t, err)

	// Verifica o status code 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Verifica se a resposta não está vazia
	assert.NotEmpty(t, Response.TempC)
	assert.NotEmpty(t, Response.TempF)
	assert.NotEmpty(t, Response.TempK)
}

// Erro 422 - "invalid zipcode"
func TestHandlerClimaCode422(t *testing.T) {
	router := chi.NewRouter()
	router.Get("/{cep}", HandlerClima)
	req, err := http.NewRequest("GET", "/340000001", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NoError(t, err)

	var errorMessage map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errorMessage)
	if err != nil {
		return
	}

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Empty(t, errorMessage)
}

// Error 404 - "can not find zipcode"
func TestHandlerClimaCode404(t *testing.T) {
	router := chi.NewRouter()
	router.Get("/{cep}", HandlerClima)
	req, err := http.NewRequest("GET", "/34000000", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NoError(t, err)

	var errorMessage map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errorMessage)
	if err != nil {
		return
	}

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Empty(t, errorMessage)
}
