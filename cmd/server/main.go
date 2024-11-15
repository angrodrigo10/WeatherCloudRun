package main

import (
	"fmt"
	"net/http"

	"github.com/angrodrigo10/WeatherCloudRun/internal/infra/web/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HandlerClima)
	fmt.Println("Servidor escutando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}
