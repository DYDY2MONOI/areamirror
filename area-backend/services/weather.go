package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type WeatherService struct {
	apiKey string
}

type WeatherResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Rain struct {
		OneH float64 `json:"1h"`
	} `json:"rain"`
	Snow struct {
		OneH float64 `json:"1h"`
	} `json:"snow"`
	Visibility int `json:"visibility"`
	Sys        struct {
		Sunrise int64 `json:"sunrise"`
		Sunset  int64 `json:"sunset"`
	} `json:"sys"`
	Name string `json:"name"`
}

type WeatherTriggerConfig struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	Operator    string  `json:"operator"`
}

type WeatherTriggerResult struct {
	Triggered bool                   `json:"triggered"`
	Data      map[string]interface{} `json:"data"`
	Message   string                 `json:"message"`
}

func NewWeatherService() (*WeatherService, error) {
	return &WeatherService{
		apiKey: "openmeteo",
	}, nil
}

func (w *WeatherService) GetCurrentWeather(city string) (*WeatherResponse, error) {
	if w.apiKey == "openmeteo" {
		return w.getOpenMeteoData(city)
	}

	return w.getMockWeatherData(city), nil
}

func (w *WeatherService) CheckWeatherTrigger(config WeatherTriggerConfig) (*WeatherTriggerResult, error) {
	weather, err := w.GetCurrentWeather(config.City)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %v", err)
	}

	result := &WeatherTriggerResult{
		Triggered: false,
		Data: map[string]interface{}{
			"city":         weather.Name,
			"temperature":  weather.Main.Temp,
			"condition":    weather.Weather[0].Main,
			"description":  weather.Weather[0].Description,
			"humidity":     weather.Main.Humidity,
			"pressure":     weather.Main.Pressure,
			"wind_speed":   weather.Wind.Speed,
			"visibility":   weather.Visibility,
			"timestamp":    time.Now().Unix(),
		},
		Message: "",
	}

	if config.Temperature > 0 {
		switch config.Operator {
		case "greater_than":
			if weather.Main.Temp > config.Temperature {
				result.Triggered = true
				result.Message = fmt.Sprintf("Temperature %.1f°C is greater than threshold %.1f°C", weather.Main.Temp, config.Temperature)
			}
		case "less_than":
			if weather.Main.Temp < config.Temperature {
				result.Triggered = true
				result.Message = fmt.Sprintf("Temperature %.1f°C is less than threshold %.1f°C", weather.Main.Temp, config.Temperature)
			}
		case "equals":
			if weather.Main.Temp == config.Temperature {
				result.Triggered = true
				result.Message = fmt.Sprintf("Temperature %.1f°C equals threshold %.1f°C", weather.Main.Temp, config.Temperature)
			}
		default:
			if weather.Main.Temp > config.Temperature {
				result.Triggered = true
				result.Message = fmt.Sprintf("Temperature %.1f°C is greater than threshold %.1f°C", weather.Main.Temp, config.Temperature)
			}
		}
	}

	if config.Condition != "" {
		conditionMatch := false
		switch config.Operator {
		case "contains":
			conditionMatch = contains(weather.Weather[0].Main, config.Condition) || contains(weather.Weather[0].Description, config.Condition)
		case "equals":
			conditionMatch = weather.Weather[0].Main == config.Condition
		default:
			conditionMatch = weather.Weather[0].Main == config.Condition
		}

		if conditionMatch {
			result.Triggered = true
			if result.Message != "" {
				result.Message += " and "
			}
			result.Message += fmt.Sprintf("Weather condition '%s' matches trigger", weather.Weather[0].Main)
		}
	}

	return result, nil
}

type OpenMeteoGeocodingResponse struct {
	Results []struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Country  string  `json:"country"`
		Latitude float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type OpenMeteoWeatherResponse struct {
	Current struct {
		Time          string  `json:"time"`
		Temperature2m float64 `json:"temperature_2m"`
		RelativeHumidity2m int `json:"relative_humidity_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		Precipitation float64 `json:"precipitation"`
		WeatherCode   int     `json:"weather_code"`
		PressureMSL   float64 `json:"pressure_msl"`
		SurfacePressure float64 `json:"surface_pressure"`
		WindSpeed10m  float64 `json:"wind_speed_10m"`
		WindDirection10m float64 `json:"wind_direction_10m"`
		Visibility    float64 `json:"visibility"`
	} `json:"current"`
	CurrentUnits struct {
		Temperature2m string `json:"temperature_2m"`
		RelativeHumidity2m string `json:"relative_humidity_2m"`
		ApparentTemperature string `json:"apparent_temperature"`
		Precipitation string `json:"precipitation"`
		PressureMSL string `json:"pressure_msl"`
		SurfacePressure string `json:"surface_pressure"`
		WindSpeed10m string `json:"wind_speed_10m"`
		WindDirection10m string `json:"wind_direction_10m"`
		Visibility string `json:"visibility"`
	} `json:"current_units"`
}

func (w *WeatherService) getOpenMeteoData(city string) (*WeatherResponse, error) {
	encodedCity := url.QueryEscape(city)
	geocodingURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json", encodedCity)

	resp, err := http.Get(geocodingURL)
	if err != nil {
		log.Printf("Failed to fetch geocoding data for %s: %v", city, err)
		return w.getMockWeatherData(city), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Open-Meteo geocoding API returned status %d for %s", resp.StatusCode, city)
		return w.getMockWeatherData(city), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read geocoding response body: %v", err)
		return w.getMockWeatherData(city), nil
	}

	var geocoding OpenMeteoGeocodingResponse
	if err := json.Unmarshal(body, &geocoding); err != nil {
		log.Printf("Failed to parse geocoding data: %v", err)
		return w.getMockWeatherData(city), nil
	}

	if len(geocoding.Results) == 0 {
		log.Printf("No location found for %s", city)
		return w.getMockWeatherData(city), nil
	}

	lat := geocoding.Results[0].Latitude
	lon := geocoding.Results[0].Longitude
	weatherURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,relative_humidity_2m,apparent_temperature,precipitation,weather_code,pressure_msl,surface_pressure,wind_speed_10m,wind_direction_10m,visibility", lat, lon)

	weatherResp, err := http.Get(weatherURL)
	if err != nil {
		log.Printf("Failed to fetch weather data for %s: %v", city, err)
		return w.getMockWeatherData(city), nil
	}
	defer weatherResp.Body.Close()

	if weatherResp.StatusCode != http.StatusOK {
		log.Printf("Open-Meteo weather API returned status %d for %s", weatherResp.StatusCode, city)
		return w.getMockWeatherData(city), nil
	}

	weatherBody, err := io.ReadAll(weatherResp.Body)
	if err != nil {
		log.Printf("Failed to read weather response body: %v", err)
		return w.getMockWeatherData(city), nil
	}

	var openMeteo OpenMeteoWeatherResponse
	if err := json.Unmarshal(weatherBody, &openMeteo); err != nil {
		log.Printf("Failed to parse weather data: %v", err)
		return w.getMockWeatherData(city), nil
	}

	condition := w.getWeatherConditionFromCode(openMeteo.Current.WeatherCode)

	return &WeatherResponse{
		Main: struct {
			Temp     float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin  float64 `json:"temp_min"`
			TempMax  float64 `json:"temp_max"`
			Pressure int     `json:"pressure"`
			Humidity int     `json:"humidity"`
		}{
			Temp:      openMeteo.Current.Temperature2m,
			FeelsLike: openMeteo.Current.ApparentTemperature,
			TempMin:   openMeteo.Current.Temperature2m - 2,
			TempMax:   openMeteo.Current.Temperature2m + 2,
			Pressure:  int(openMeteo.Current.PressureMSL),
			Humidity:  openMeteo.Current.RelativeHumidity2m,
		},
		Weather: []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		}{
			{
				ID:          openMeteo.Current.WeatherCode,
				Main:        condition,
				Description: condition,
				Icon:        fmt.Sprintf("%02d", openMeteo.Current.WeatherCode),
			},
		},
		Wind: struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		}{
			Speed: openMeteo.Current.WindSpeed10m,
			Deg:   int(openMeteo.Current.WindDirection10m),
		},
		Clouds: struct {
			All int `json:"all"`
		}{
			All: 0,
		},
		Rain: struct {
			OneH float64 `json:"1h"`
		}{
			OneH: openMeteo.Current.Precipitation,
		},
		Snow: struct {
			OneH float64 `json:"1h"`
		}{
			OneH: 0,
		},
		Visibility: int(openMeteo.Current.Visibility * 1000),
		Sys: struct {
			Sunrise int64 `json:"sunrise"`
			Sunset  int64 `json:"sunset"`
		}{
			Sunrise: time.Now().Unix() - 3600,
			Sunset:  time.Now().Unix() + 3600,
		},
		Name: geocoding.Results[0].Name,
	}, nil
}

func (w *WeatherService) getWeatherConditionFromCode(code int) string {
	switch {
	case code == 0:
		return "Clear"
	case code == 1, code == 2, code == 3:
		return "Clouds"
	case code == 45, code == 48:
		return "Fog"
	case code == 51, code == 53, code == 55:
		return "Drizzle"
	case code == 56, code == 57:
		return "Drizzle"
	case code == 61, code == 63, code == 65:
		return "Rain"
	case code == 66, code == 67:
		return "Rain"
	case code == 71, code == 73, code == 75:
		return "Snow"
	case code == 77:
		return "Snow"
	case code == 80, code == 81, code == 82:
		return "Rain"
	case code == 85, code == 86:
		return "Snow"
	case code == 95:
		return "Thunderstorm"
	case code == 96, code == 99:
		return "Thunderstorm"
	default:
		return "Clear"
	}
}

type MetaWeatherLocation struct {
	Woeid int `json:"woeid"`
	Title string `json:"title"`
}

type MetaWeatherResponse struct {
	ConsolidatedWeather []struct {
		ID                   int64   `json:"id"`
		WeatherStateName     string  `json:"weather_state_name"`
		WeatherStateAbbr     string  `json:"weather_state_abbr"`
		WindDirection        float64 `json:"wind_direction"`
		WindDirectionCompass string  `json:"wind_direction_compass"`
		WindSpeed            float64 `json:"wind_speed"`
		TheTemp              float64 `json:"the_temp"`
		MinTemp              float64 `json:"min_temp"`
		MaxTemp              float64 `json:"max_temp"`
		AirPressure          float64 `json:"air_pressure"`
		Humidity             int     `json:"humidity"`
		Visibility           float64 `json:"visibility"`
		Predictability       int     `json:"predictability"`
	} `json:"consolidated_weather"`
	Title string `json:"title"`
}

func (w *WeatherService) getMetaWeatherData(city string) (*WeatherResponse, error) {
	locationURL := fmt.Sprintf("https://www.metaweather.com/api/location/search/?query=%s", city)

	resp, err := http.Get(locationURL)
	if err != nil {
		log.Printf("Failed to fetch location data for %s: %v", city, err)
		return w.getMockWeatherData(city), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("MetaWeather location API returned status %d for %s", resp.StatusCode, city)
		return w.getMockWeatherData(city), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read location response body: %v", err)
		return w.getMockWeatherData(city), nil
	}

	var locations []MetaWeatherLocation
	if err := json.Unmarshal(body, &locations); err != nil {
		log.Printf("Failed to parse location data: %v", err)
		return w.getMockWeatherData(city), nil
	}

	if len(locations) == 0 {
		log.Printf("No location found for %s", city)
		return w.getMockWeatherData(city), nil
	}

	woeid := locations[0].Woeid
	weatherURL := fmt.Sprintf("https://www.metaweather.com/api/location/%d/", woeid)

	weatherResp, err := http.Get(weatherURL)
	if err != nil {
		log.Printf("Failed to fetch weather data for woeid %d: %v", woeid, err)
		return w.getMockWeatherData(city), nil
	}
	defer weatherResp.Body.Close()

	if weatherResp.StatusCode != http.StatusOK {
		log.Printf("MetaWeather weather API returned status %d for woeid %d", weatherResp.StatusCode, woeid)
		return w.getMockWeatherData(city), nil
	}

	weatherBody, err := io.ReadAll(weatherResp.Body)
	if err != nil {
		log.Printf("Failed to read weather response body: %v", err)
		return w.getMockWeatherData(city), nil
	}

	var metaWeather MetaWeatherResponse
	if err := json.Unmarshal(weatherBody, &metaWeather); err != nil {
		log.Printf("Failed to parse weather data: %v", err)
		return w.getMockWeatherData(city), nil
	}

	if len(metaWeather.ConsolidatedWeather) == 0 {
		log.Printf("No weather data found for woeid %d", woeid)
		return w.getMockWeatherData(city), nil
	}

	current := metaWeather.ConsolidatedWeather[0]

	return &WeatherResponse{
		Main: struct {
			Temp     float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin  float64 `json:"temp_min"`
			TempMax  float64 `json:"temp_max"`
			Pressure int     `json:"pressure"`
			Humidity int     `json:"humidity"`
		}{
			Temp:      current.TheTemp,
			FeelsLike: current.TheTemp,
			TempMin:   current.MinTemp,
			TempMax:   current.MaxTemp,
			Pressure:  int(current.AirPressure),
			Humidity:  current.Humidity,
		},
		Weather: []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		}{
			{
				ID:          int(current.ID),
				Main:        current.WeatherStateName,
				Description: current.WeatherStateName,
				Icon:        current.WeatherStateAbbr,
			},
		},
		Wind: struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		}{
			Speed: current.WindSpeed,
			Deg:   int(current.WindDirection),
		},
		Clouds: struct {
			All int `json:"all"`
		}{
			All: 0,
		},
		Rain: struct {
			OneH float64 `json:"1h"`
		}{
			OneH: 0,
		},
		Snow: struct {
			OneH float64 `json:"1h"`
		}{
			OneH: 0,
		},
		Visibility: int(current.Visibility * 1000),
		Sys: struct {
			Sunrise int64 `json:"sunrise"`
			Sunset  int64 `json:"sunset"`
		}{
			Sunrise: time.Now().Unix() - 3600,
			Sunset:  time.Now().Unix() + 3600,
		},
		Name: metaWeather.Title,
	}, nil
}

func (w *WeatherService) getMockWeatherData(city string) *WeatherResponse {
	return &WeatherResponse{
		Main: struct {
			Temp     float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin  float64 `json:"temp_min"`
			TempMax  float64 `json:"temp_max"`
			Pressure int     `json:"pressure"`
			Humidity int     `json:"humidity"`
		}{
			Temp:      17.0,
			FeelsLike: 16.5,
			TempMin:   15.0,
			TempMax:   19.0,
			Pressure:  1013,
			Humidity:  75,
		},
		Weather: []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		}{
			{
				ID:          500,
				Main:        "Rain",
				Description: "light rain",
				Icon:        "10d",
			},
		},
		Wind: struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		}{
			Speed: 4.2,
			Deg:   220,
		},
		Clouds: struct {
			All int `json:"all"`
		}{
			All: 80,
		},
		Rain: struct {
			OneH float64 `json:"1h"`
		}{
			OneH: 0.5,
		},
		Snow: struct {
			OneH float64 `json:"1h"`
		}{
			OneH: 0,
		},
		Visibility: 10000,
		Sys: struct {
			Sunrise int64 `json:"sunrise"`
			Sunset  int64 `json:"sunset"`
		}{
			Sunrise: time.Now().Unix() - 3600,
			Sunset:  time.Now().Unix() + 3600,
		},
		Name: city,
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (w *WeatherService) TestWeatherTrigger(config WeatherTriggerConfig) (*WeatherTriggerResult, error) {
	log.Printf("Testing weather trigger for city: %s, temperature: %.1f, condition: %s",
		config.City, config.Temperature, config.Condition)

	return w.CheckWeatherTrigger(config)
}
