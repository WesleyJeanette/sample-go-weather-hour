package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ShortForecast struct {
	Latitude        float64
	Longitude       float64
	Name            string
	StartTime       string
	EndTime         string
	Temperature     int
	TemperatureUnit string
	ShortForecast   string
}

// ForecastURLs (I don't care for this name)
type ForecastURLS struct {
	City      string
	State     string
	Point     string // the point URL for the forecast
	Latitude  float64
	Longitude float64

	// URLs for different forecast types
	Daily    string
	Hourly   string
	GridData string
}

type WeatherService interface {
	// RequestWeatherData fetches weather data for the given latitude and longitude.
	RequestWeatherData(ctx context.Context, lat float64, lon float64) (string, error)
}

// NationalWeatherService implements the WeatherService interface
type NationalWeatherService struct {
	apiURL string
	apiKey string
}

func NewNationalWeatherService() NationalWeatherService {
	return NationalWeatherService{
		apiURL: "https://api.weather.gov/",
	}
}

// RequestWeatherData
func (ws NationalWeatherService) RequestWeatherData(ctx context.Context, lat, lon float64) (ForecastURLS, error) {
	forecasts := ForecastURLS{
		Latitude:  lat, // default to input lat
		Longitude: lon, // default to input lon
	}
	// The API endpoint for points will return a set of properties including a forecast URL for the
	// location closest to the given latitude and longitude.
	pointsURL := fmt.Sprintf("%spoints/%f,%f", ws.apiURL, lat, lon)
	resp, err := http.Get(pointsURL)
	if err != nil {
		return forecasts, fmt.Errorf("failed to fetch weather data: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return forecasts, fmt.Errorf("failed to fetch weather data: received status code %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return forecasts, fmt.Errorf("failed to read response body: %v", err)
	}
	// At this point, for a quick project, all that is needed is to extract the forecast URL from the points response.
	type PointsResponse struct {
		Properties struct {
			Daily            string `json:"forecast"`
			Hourly           string `json:"forecastHourly"`
			GridData         string `json:"forecastGridData"`
			RelativeLocation struct {
				Properties struct {
					City  string `json:"city"`
					State string `json:"state"`
				} `json:"properties"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"relativeLocation"`
		} `json:"properties"`
	}
	var pointsResp PointsResponse
	err = json.Unmarshal(data, &pointsResp)
	if err != nil {
		return forecasts, fmt.Errorf("failed to parse points response: %v", err)
	}

	fmt.Printf("Points Response: %+v\n", pointsResp)
	// Populate the ForecastURLS struct with the relevant data
	forecasts.City = pointsResp.Properties.RelativeLocation.Properties.City
	forecasts.State = pointsResp.Properties.RelativeLocation.Properties.State
	if len(pointsResp.Properties.RelativeLocation.Coordinates) == 2 {
		forecasts.Longitude = pointsResp.Properties.RelativeLocation.Coordinates[0] // lon is first in the array
		forecasts.Latitude = pointsResp.Properties.RelativeLocation.Coordinates[1]  // lat is second in the array
	}

	forecasts.Point = pointsURL
	forecasts.Daily = pointsResp.Properties.Daily
	forecasts.Hourly = pointsResp.Properties.Hourly
	forecasts.GridData = pointsResp.Properties.GridData

	return forecasts, nil
}

func (ws NationalWeatherService) RequestShortForcastData(ctx context.Context, fURLs ForecastURLS) (ShortForecast, error) {
	sf := ShortForecast{
		Latitude:  fURLs.Latitude,
		Longitude: fURLs.Longitude,
	}
	// Use the Daily forecast URL for this, it would be neat if we passed in what type of forecast we want
	forecastResp, err := http.Get(fURLs.Daily)
	if err != nil {
		return sf, fmt.Errorf("failed to fetch forecast data: %v", err)
	}
	defer forecastResp.Body.Close()
	if forecastResp.StatusCode != http.StatusOK {
		return sf, fmt.Errorf("failed to fetch forecast data: received status code %d", forecastResp.StatusCode)
	}
	data, err := io.ReadAll(forecastResp.Body)
	if err != nil {
		return sf, fmt.Errorf("failed to read forecast response body: %v", err)
	}
	type forecastResponse struct {
		Properties struct {
			Periods []struct {
				Name            string `json:"name"`
				StartTime       string `json:"startTime"`
				EndTime         string `json:"endTime"`
				Temperature     int    `json:"temperature"`
				TemperatureUnit string `json:"temperatureUnit"`
				ShortForecast   string `json:"shortForecast"`
			} `json:"periods"`
		} `json:"properties"`
	}
	var forecastRespData forecastResponse
	err = json.Unmarshal(data, &forecastRespData)
	if err != nil {
		return sf, fmt.Errorf("failed to parse forecast response: %v", err)
	}

	// for this quick project, return the first forecast period
	if len(forecastRespData.Properties.Periods) > 0 {
		sf.Name = forecastRespData.Properties.Periods[0].Name
		sf.StartTime = forecastRespData.Properties.Periods[0].StartTime
		sf.EndTime = forecastRespData.Properties.Periods[0].EndTime
		sf.Temperature = forecastRespData.Properties.Periods[0].Temperature
		sf.TemperatureUnit = forecastRespData.Properties.Periods[0].TemperatureUnit
		sf.ShortForecast = forecastRespData.Properties.Periods[0].ShortForecast

		return sf, nil
	}

	return sf, fmt.Errorf("no forecast periods available")
}
