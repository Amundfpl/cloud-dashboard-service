// Package services provides logic for enriching dashboards with live and cached data.
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/amundfpl/Assignment-2/cache"
	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/httpclient"
	"github.com/amundfpl/Assignment-2/utils"
)

// GetEnrichedDashboards fetches and enriches all dashboard configs with live data.
// It adds country, weather, and currency info using cache where available.
func GetEnrichedDashboards() ([]utils.DashboardResponse, error) {
	ctx := context.Background()

	configs, configFetchErr := db.GetAllDashboardConfigs(ctx)
	if configFetchErr != nil {
		return nil, fmt.Errorf("%s: %w", utils.ErrFetchAllConfigs, configFetchErr)
	}

	client := httpclient.NewClient()
	var results []utils.DashboardResponse

	// Loop through each dashboard config and enrich with external data.
	for _, cfg := range configs {
		resp := utils.DashboardResponse{
			Country: cfg.Country,
			ISOCode: cfg.ISOCode,
		}

		// Step 1: Enrich with country-level info (capital, coords, etc.)
		countryInfo, enrichCountryErr := enrichCountryData(client, cfg, &resp)
		if enrichCountryErr != nil {
			return nil, fmt.Errorf("%s: %w", utils.ErrEnrichCountry, enrichCountryErr)
		}

		// Step 2: Enrich with weather info (temperature, precipitation)
		enrichWeatherErr := enrichWeatherData(client, cfg, &resp)
		if enrichWeatherErr != nil {
			return nil, fmt.Errorf("%s: %w", utils.ErrEnrichWeather, enrichWeatherErr)
		}

		// Step 3: Enrich with currency exchange data
		enrichCurrencyErr := enrichCurrencyData(client, cfg, countryInfo, &resp)
		if enrichCurrencyErr != nil {
			return nil, fmt.Errorf("%s: %w", utils.ErrEnrichCurrency, enrichCurrencyErr)
		}

		results = append(results, resp)
	}

	return results, nil
}

// enrichCountryData enriches a dashboard with capital, coordinates, population, and area info.
// Attempts cache first, otherwise fetches from external API and stores to cache.
func enrichCountryData(client *httpclient.Client, cfg utils.DashboardConfig, resp *utils.DashboardResponse) (utils.CountryInfoResponse, error) {
	if !(cfg.Features.Capital || cfg.Features.Coordinates || cfg.Features.Population || cfg.Features.Area) {
		return utils.CountryInfoResponse{}, nil // Nothing to enrich
	}

	cached, cacheHitErr := cache.GetCachedCountryInfo(context.Background(), cfg.ISOCode, 24*time.Hour)
	if cacheHitErr == nil {
		syncCountryFields(cfg, resp, *cached)
		return *cached, nil
	}

	countryInfo, countryFetchErr := fetchCountryInfo(client, cfg.ISOCode)
	if countryFetchErr != nil {
		return utils.CountryInfoResponse{}, countryFetchErr
	}

	_ = cache.SaveCountryInfoToCache(context.Background(), cfg.ISOCode, countryInfo)
	syncCountryFields(cfg, resp, countryInfo)
	return countryInfo, nil
}

// syncCountryFields maps selected fields from the country API response to the dashboard response struct.
func syncCountryFields(cfg utils.DashboardConfig, resp *utils.DashboardResponse, info utils.CountryInfoResponse) {
	if cfg.Features.Capital && len(info.Capital) > 0 {
		resp.Capital = info.Capital[0]
	}
	if cfg.Features.Coordinates && len(info.Latlng) == 2 {
		resp.Latitude = info.Latlng[0]
		resp.Longitude = info.Latlng[1]
	}
	if cfg.Features.Population {
		resp.Population = info.Population
	}
	if cfg.Features.Area {
		resp.Area = info.Area
	}
}

// enrichWeatherData adds temperature and precipitation values using cache or a fresh API call.
func enrichWeatherData(client *httpclient.Client, cfg utils.DashboardConfig, resp *utils.DashboardResponse) error {
	if !(cfg.Features.Temperature || cfg.Features.Precipitation) {
		return nil // Nothing to enrich
	}

	key := cache.WeatherCacheKey(resp.Latitude, resp.Longitude)
	cached, cacheErr := cache.GetCachedWeather(context.Background(), key, 2*time.Hour)
	if cacheErr == nil {
		if cfg.Features.Temperature {
			resp.Temperature = cached.Temperature
		}
		if cfg.Features.Precipitation {
			resp.Precipitation = cached.Precipitation
		}
		return nil
	}

	weather, weatherFetchErr := fetchWeather(client, resp.Latitude, resp.Longitude)
	if weatherFetchErr != nil {
		return weatherFetchErr
	}

	_ = cache.SaveWeatherToCache(context.Background(), key, weather)
	if cfg.Features.Temperature {
		resp.Temperature = weather.Temperature
	}
	if cfg.Features.Precipitation {
		resp.Precipitation = weather.Precipitation
	}
	return nil
}

// enrichCurrencyData attaches exchange rate information to a dashboard response.
func enrichCurrencyData(client *httpclient.Client, cfg utils.DashboardConfig, countryInfo utils.CountryInfoResponse, resp *utils.DashboardResponse) error {
	if len(cfg.Features.TargetCurrencies) == 0 {
		return nil // Nothing to enrich
	}

	var base string
	for currency := range countryInfo.Currencies {
		base = currency
		break
	}
	if base == "" {
		return fmt.Errorf(utils.ErrNoBaseCurrency)
	}

	ctx := context.Background()
	key := cache.CurrencyCacheKey(base, cfg.Features.TargetCurrencies)

	cached, cacheErr := cache.GetCachedCurrencyRates(ctx, key, 12*time.Hour)
	if cacheErr == nil {
		resp.ExchangeRates = cached
		return nil
	}

	rates, currencyFetchErr := fetchCurrencyRates(client, countryInfo.Currencies, cfg.Features.TargetCurrencies)
	if currencyFetchErr != nil {
		return currencyFetchErr
	}

	_ = cache.SaveCurrencyRatesToCache(ctx, key, rates)
	resp.ExchangeRates = rates
	return nil
}
