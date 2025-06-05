package handlers

import (
	"encoding/json"
	"net/http"
	"sample-go-weather-hour/internal/services"
	"strconv"
)

type Handler struct {
	NWS services.NationalWeatherService
}

func NewHandler() *Handler {
	return &Handler{
		NWS: services.NewNationalWeatherService(),
	}
}

func (h *Handler) GetForecastHandler(w http.ResponseWriter, r *http.Request) {

	// it would be nice to have a rate limiter, we are only allowed x requests per minute to the external api...
	// it would also be nice to have a cache for the results, so we don't have to hit the external api for every request

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for required query parameters: lat and lon
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")
	if lat == "" || lon == "" {
		http.Error(w, "Missing latitude or longitude parameters", http.StatusBadRequest)
		return
	}
	// Validate latitude and longitude
	// Latitude must be between -90 and 90
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil || latFloat < -90 || latFloat > 90 {
		http.Error(w, "Invalid latitude value", http.StatusBadRequest)
		return
	}
	// Longitude must be between -180 and 180
	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil || lonFloat < -180 || lonFloat > 180 {
		http.Error(w, "Invalid longitude value", http.StatusBadRequest)
		return
	}

	// If we wanted to use different weather services, we also provide a query parameter for that
	//service := r.URL.Query().Get("service")
	// but for now we only support the National Weather Service

	forcastURL, err := h.NWS.RequestWeatherData(r.Context(), latFloat, lonFloat)
	if err != nil {
		http.Error(w, "Failed to fetch weather data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// We have the forecast URL, now we can fetch the forecast data
	short, err := h.NWS.RequestShortForcastData(r.Context(), forcastURL)
	if err != nil {
		http.Error(w, "Failed to forecast data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type ForcastReport struct {
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		Name          string  `json:"name"`
		ShortForecast string  `json:"short_forecast"`
		TempVibe      string  `json:"temp_vibe"`
	}

	report := ForcastReport{
		Latitude:      latFloat,
		Longitude:     lonFloat,
		Name:          short.Name,
		ShortForecast: short.ShortForecast,
	}
	if short.TemperatureUnit == "C" {
		report.TempVibe = "Convert to Fahrenheit for vibe" // Placeholder for conversion logic
	} else if short.TemperatureUnit == "F" {
		if short.Temperature < 32 {
			report.TempVibe = "freezing"
		} else if short.Temperature < 50 {
			report.TempVibe = "cold"
		} else if short.Temperature < 70 {
			report.TempVibe = "perfect"
		} else if short.Temperature < 80 {
			report.TempVibe = "warm"
		} else {
			report.TempVibe = "hot"
		}
	} else {
		report.TempVibe = "unknown" // handle unexpected temperature unit
	}
	// Implement the logic to get the weather forecast
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(report)
	if err != nil {
		http.Error(w, "Failed to marshal forecast report: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (h *Handler) GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to get the service status
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service is running"))
}
