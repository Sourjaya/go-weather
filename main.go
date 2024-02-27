package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type apiConfigData struct {
	RapidApiKey string `json:"RapidApiKey"`
}
type weatherData struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		Celcius   float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}
	var c apiConfigData
	if err := json.Unmarshal(bytes, &c); err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}
func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}

	url := "https://weatherapi-com.p.rapidapi.com/current.json?q=" + city

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", apiConfig.RapidApiKey)
	req.Header.Add("X-RapidAPI-Host", "weatherapi-com.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return weatherData{}, err
	}

	defer res.Body.Close()
	var d weatherData
	if err := json.NewDecoder(res.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello\n"))
}
func weather(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city := vars["city"]
	data, err := query(city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(data)
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello)
	r.HandleFunc("/weather/{city}", weather)
	fmt.Printf("Starting server at port 8080\n")
	http.ListenAndServe(":8080", r)
}
