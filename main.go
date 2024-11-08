package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type WeatherResponse struct {
	TemperatureC float64 `json:"temp_c"`
}

type CepResponse struct {
	Localidade string `json:"localidade"`
}

func main() {
	http.HandleFunc("/clima", handleClima)
	fmt.Println("Servidor escutando na porta 8081...")
	http.ListenAndServe(":8081", nil)
}

func handleClima(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")

	if !isValidCep(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	localidade, err := getLocalidade(cep)
	fmt.Printf("Localidade: %v\n", localidade)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	temperature, err := getTemperature(localidade)
	fmt.Printf("Localidade: %v\n", err)
	if err != nil {
		http.Error(w, "could not fetch weather", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Temperatura em %v\n", temperature)

	response := map[string]float64{
		"temp_C": temperature,
		"temp_F": celsiusToFahrenheit(temperature),
		"temp_K": celsiusToKelvin(temperature),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func isValidCep(cep string) bool {
	re := regexp.MustCompile(`^\d{5}-?\d{3}$`)
	return re.MatchString(cep)
}

func getLocalidade(cep string) (string, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid cep")
	}

	var cepResp CepResponse
	if err := json.NewDecoder(resp.Body).Decode(&cepResp); err != nil {
		return "", err
	}

	return cepResp.Localidade, nil
}

func getTemperature(city string) (float64, error) {
	apiKey := "1da14e1d67344108ab9194641242010" // Insira sua chave da API WeatherAPI aqui
	//resp, err := http.Get(fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city))
	resp, err := http.Get("http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + city)

	//url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city)
	//resp, err := http.Get(url)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	fmt.Println("http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + city)
	fmt.Printf("\nErro http %v\n", resp.StatusCode)
	fmt.Printf("\nresp %v\n", resp.Body)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("could not fetch weather data")
	}

	var weatherResp WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return 0, err
	}

	return weatherResp.TemperatureC, nil
}

func celsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func celsiusToKelvin(c float64) float64 {
	return c + 273.15
}
