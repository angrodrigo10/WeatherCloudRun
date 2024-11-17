package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/angrodrigo10/WeatherCloudRun/internal/infra/web/handlers"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Get("/{cep}", handlers.HandlerClima)
	log.Println("Servidor iniciado na porta 8080!")
	http.ListenAndServe(":8080", router)

	//mux := http.NewServeMux()
	//mux.HandleFunc("GET /{cep}", handlers.HandlerClima)
	//fmt.Println("Servidor escutando na porta 8080...")
	//http.ListenAndServe(":8080", mux)

	//http.HandleFunc("/", handlers.HandlerClima)
	//fmt.Println("Servidor escutando na porta 8080...")
	//http.ListenAndServe(":8080", nil)
}
