package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/go-chi/chi"
)

type CurrentWeather struct {
	TemperatureC float64 `json:"temp_c"`
}

type WeatherResponse struct {
	Current CurrentWeather `json:"current"`
}

type CepResponse struct {
	Localidade string `json:"localidade"`
}

func HandlerClima(w http.ResponseWriter, r *http.Request) {
	//cep := r.URL.Query().Get("cep")
	//cep := r.PathValue("cep")
	cep := chi.URLParam(r, "cep")

	fmt.Printf("CEP: %v\n", cep)

	if !isValidCep(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	localidade, err := getLocalidade(cep)

	fmt.Printf("localidade: %v\n", localidade)
	fmt.Printf("err: %v\n", err)

	if err != nil || localidade == "" {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	encodedLocalidade := url.QueryEscape(localidade)

	fmt.Printf("Localidade: %v\n", encodedLocalidade)

	temperature, err := getTemperature(encodedLocalidade)
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
	apiKey := "1da14e1d67344108ab9194641242010"
	resp, err := http.Get("http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + city)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("could not fetch weather data")
	}

	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		return 0, err
	}

	return weatherResp.Current.TemperatureC, nil
}

func celsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func celsiusToKelvin(c float64) float64 {
	return c + 273.15
}
