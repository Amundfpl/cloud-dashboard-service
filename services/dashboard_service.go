// Package services provides the core business logic for dashboard operations.
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/httpclient"
	"github.com/amundfpl/Assignment-2/utils"
)

// DashboardService defines an interface for dashboard operations.
type DashboardService interface {
	GetPopulatedDashboardByID(id string) (*utils.PopulatedDashboardResponse, error)
	GetEnrichedDashboards() ([]utils.DashboardResponse, error)
}

// RealDashboardService is a concrete implementation of DashboardService.
type RealDashboardService struct{}

func (r RealDashboardService) GetPopulatedDashboardByID(id string) (*utils.PopulatedDashboardResponse, error) {
	return GetPopulatedDashboardByID(id)
}

func (r RealDashboardService) GetEnrichedDashboards() ([]utils.DashboardResponse, error) {
	return GetEnrichedDashboards()
}

// GetPopulatedDashboardByID builds a full dashboard response by enriching a config with live data.
func GetPopulatedDashboardByID(id string) (*utils.PopulatedDashboardResponse, error) {
	ctx := context.Background()

	// Step 1: Retrieve dashboard config from Firestore
	config, fetchErr := db.GetDashboardConfigByID(ctx, id)
	if fetchErr != nil {
		return nil, fmt.Errorf("%s: %w", utils.ErrFetchConfig, fetchErr)
	}

	// Step 2: Initialize response base
	client := httpclient.NewClient()
	resp := &utils.PopulatedDashboardResponse{
		Country: config.Country,
		ISOCode: config.ISOCode,
	}

	// Step 3: Fetch country information from API
	countryInfo, countryErr := fetchCountryInfo(client, config.ISOCode)
	if countryErr != nil {
		return nil, fmt.Errorf("%s: %w", utils.ErrInvalidCountryResp, countryErr)
	}

	features := utils.PopulatedFeatures{}

	// Step 4: Populate capital if enabled
	if config.Features.Capital && len(countryInfo.Capital) > 0 {
		features.Capital = countryInfo.Capital[0]
	}

	// Step 5: Populate coordinates if enabled
	if config.Features.Coordinates && len(countryInfo.Latlng) == 2 {
		features.Coordinates = &utils.Coordinates{
			Latitude:  countryInfo.Latlng[0],
			Longitude: countryInfo.Latlng[1],
		}
	}

	// Step 6: Add population if requested
	if config.Features.Population {
		features.Population = countryInfo.Population
	}

	// Step 7: Add area if requested
	if config.Features.Area {
		features.Area = countryInfo.Area
	}

	// Step 8: If weather data is needed and coordinates are available, fetch it
	if (config.Features.Temperature || config.Features.Precipitation) && features.Coordinates != nil {
		weather, weatherErr := fetchWeather(client, features.Coordinates.Latitude, features.Coordinates.Longitude)
		if weatherErr != nil {
			return nil, fmt.Errorf("%s: %w", utils.ErrFetchWeather, weatherErr)
		}

		// Step 9: Add temperature if requested, trigger LOW_TEMP webhook if under 0°C
		if config.Features.Temperature {
			features.Temperature = weather.Temperature
			if weather.Temperature < 0 {
				TriggerWebhooks(utils.EventLowTemp, config.ISOCode)
			}
		}

		// Step 10: Add precipitation if requested
		if config.Features.Precipitation {
			features.Precipitation = weather.Precipitation
		}
	}

	// Step 11: If currency data is requested, fetch exchange rates
	if len(config.Features.TargetCurrencies) > 0 {
		rates, currencyErr := fetchCurrencyRates(client, countryInfo.Currencies, config.Features.TargetCurrencies)
		if currencyErr != nil {
			return nil, fmt.Errorf("%s: %w", utils.ErrFetchCurrency, currencyErr)
		}
		features.TargetCurrencies = rates
	}

	// Step 12: Finalize response
	resp.Features = features
	resp.LastRetrieval = utils.CurrentTimestamp()

	// Step 13: Trigger INVOKE webhook for dashboard access
	TriggerWebhooks(utils.EventInvoke, config.ISOCode)
	return resp, nil
}

// fetchCountryInfo retrieves country metadata from the external REST Countries API.
func fetchCountryInfo(client *httpclient.Client, isoCode string) (utils.CountryInfoResponse, error) {
	url := utils.RESTCountriesAPI + utils.RESTCountriesByAlpha + strings.ToUpper(isoCode)
	body, getErr := client.Get(url)
	if getErr != nil {
		return utils.CountryInfoResponse{}, fmt.Errorf("%s: %w", utils.ErrFetchCountry, getErr)
	}

	var countries []utils.CountryInfoResponse
	if decodeErr := json.Unmarshal(body, &countries); decodeErr != nil || len(countries) == 0 {
		return utils.CountryInfoResponse{}, fmt.Errorf("%s: %w", utils.ErrInvalidCountryResp, decodeErr)
	}

	return countries[0], nil
}

// fetchWeather retrieves current weather metrics (temperature and precipitation) for the provided coordinates.
func fetchWeather(client *httpclient.Client, lat, lon float64) (utils.WeatherData, error) {
	url := fmt.Sprintf(utils.OpenMeteoWeatherURLFmt, utils.OpenMeteoAPI, utils.OpenMeteoForecast, lat, lon)
	body, weatherErr := client.Get(url)
	if weatherErr != nil {
		return utils.WeatherData{}, fmt.Errorf("%s: %w", utils.ErrFetchWeather, weatherErr)
	}

	var result struct {
		Current struct {
			Temperature   float64 `json:"temperature_2m"`
			Precipitation float64 `json:"precipitation"`
		} `json:"current"`
	}

	if decodeErr := json.Unmarshal(body, &result); decodeErr != nil {
		return utils.WeatherData{}, fmt.Errorf("%s: %w", utils.ErrInvalidWeatherResp, decodeErr)
	}

	return utils.WeatherData{
		Temperature:   result.Current.Temperature,
		Precipitation: result.Current.Precipitation,
	}, nil
}

// fetchCurrencyRates retrieves and filters currency exchange rates based on a base and target currencies.
func fetchCurrencyRates(client *httpclient.Client, baseCurrencies map[string]utils.CurrencyDetails, targets []string) (map[string]float64, error) {
	var base string
	for k := range baseCurrencies {
		base = k
		break
	}
	if base == "" {
		return nil, fmt.Errorf(utils.ErrNoBaseCurrency)
	}

	targetStr := strings.Join(targets, ",")
	url := fmt.Sprintf("https://api.frankfurter.app/latest?from=%s&to=%s", base, targetStr)

	body, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", utils.ErrFetchCurrency, err)
	}

	var response struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("%s: %w", utils.ErrInvalidCurrencyResp, err)
	}

	return response.Rates, nil
}
