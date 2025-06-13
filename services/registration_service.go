package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/httpclient"
	"github.com/amundfpl/Assignment-2/utils"
)

// RegisterDashboardConfig processes a registration payload, resolves the country name if needed,
// stores the config in Firestore, and triggers a webhook for the REGISTER event.
func RegisterDashboardConfig(payload []byte) (map[string]string, error) {
	var request utils.RegistrationRequest

	if err := json.Unmarshal(payload, &request); err != nil {
		return nil, fmt.Errorf(utils.ErrInvalidJSONFormat, err)
	}

	// Ensure either Country or ISOCode is provided
	if request.Country == "" && request.ISOCode == "" {
		return nil, errors.New(utils.ErrMissingCountryOrISOCode)
	}

	client := httpclient.NewClient()
	countryName := request.Country

	// Resolve country name from ISOCode if needed
	if countryName == "" {
		resolvedName, err := getCountryNameByISO(client, request.ISOCode)
		if err != nil {
			return nil, err
		}
		countryName = resolvedName
	}

	// Construct the dashboard configuration
	config := utils.DashboardConfig{
		Country:    countryName,
		ISOCode:    request.ISOCode,
		Features:   request.Features,
		LastChange: time.Now().Format(utils.TimestampLayout),
	}

	// Store config in Firestore
	id, err := db.SaveDashboardConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrFirestoreSaveFailed, err)
	}

	// Trigger webhook for registration event
	TriggerWebhooks(utils.EventRegister, config.ISOCode)

	return map[string]string{
		utils.KeyID:         id,
		utils.KeyLastChange: config.LastChange,
	}, nil
}

// getCountryNameByISO queries the REST Countries API using an ISO code and returns the full country name.
func getCountryNameByISO(client *httpclient.Client, isoCode string) (string, error) {
	url := fmt.Sprintf("%s/alpha/%s", utils.RESTCountriesAPI, isoCode)

	responseData, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf(utils.ErrRESTCountryFetchFailed, err)
	}

	var response []map[string]interface{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return "", fmt.Errorf(utils.ErrCountryResponseParseFailed, err)
	}

	if len(response) == 0 || response[0]["name"] == nil {
		return "", fmt.Errorf(utils.ErrInvalidISOCode, isoCode)
	}

	name := response[0]["name"].(map[string]interface{})["common"].(string)
	return name, nil
}

// UpdateDashboardConfig replaces the entire dashboard configuration with the provided update.
// Triggers a CHANGE webhook upon success.
func UpdateDashboardConfig(ctx context.Context, id string, body []byte) (map[string]string, error) {
	var updatedConfig utils.DashboardConfig

	if err := json.Unmarshal(body, &updatedConfig); err != nil {
		return nil, fmt.Errorf(utils.ErrInvalidJSONBodyFormat, err)
	}

	updatedConfig.ID = id
	updatedConfig.LastChange = time.Now().Format(utils.TimestampLayout)

	if err := db.UpdateDashboardConfig(ctx, updatedConfig); err != nil {
		return nil, fmt.Errorf(utils.ErrFirestoreUpdateFailed, err)
	}

	TriggerWebhooks(utils.EventChange, updatedConfig.ISOCode)

	return map[string]string{
		utils.KeyID:         id,
		utils.KeyLastChange: updatedConfig.LastChange,
	}, nil
}

// PatchDashboardConfig applies a partial update to an existing dashboard configuration.
// It allows updating the country, ISO code, and individual feature flags. Triggers a PATCH webhook.
func PatchDashboardConfig(ctx context.Context, id string, patch map[string]interface{}) (map[string]string, error) {
	existingConfig, err := db.GetDashboardConfigByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrConfigNotFoundByID, err)
	}

	// Update top-level fields if present in patch
	if country, ok := patch[utils.KeyCountry].(string); ok {
		existingConfig.Country = country
	}
	if isoCode, ok := patch[utils.KeyISOCode].(string); ok {
		existingConfig.ISOCode = isoCode
	}

	// Apply patch to nested feature configuration
	if features, ok := patch[utils.KeyFeatures].(map[string]interface{}); ok {
		applyFeaturePatch(&existingConfig.Features, features)
	}

	existingConfig.LastChange = time.Now().Format(utils.TimestampLayout)

	if err := db.UpdateDashboardConfig(ctx, *existingConfig); err != nil {
		return nil, fmt.Errorf(utils.ErrFirestoreUpdateFailed, err)
	}

	TriggerWebhooks(utils.EventPatch, existingConfig.ISOCode)

	return map[string]string{
		utils.KeyID:         existingConfig.ID,
		utils.KeyLastChange: existingConfig.LastChange,
	}, nil
}

// applyFeaturePatch updates fields in FeatureConfig based on the incoming patch data.
func applyFeaturePatch(dest *utils.FeatureConfig, patch map[string]interface{}) {
	log.Println("applyFeaturePatch - input:", patch)

	if v, ok := patch[utils.KeyTemperature].(bool); ok {
		dest.Temperature = v
	}
	if v, ok := patch[utils.KeyPrecipitation].(bool); ok {
		dest.Precipitation = v
	}
	if v, ok := patch[utils.KeyCapital].(bool); ok {
		dest.Capital = v
	}
	if v, ok := patch[utils.KeyCoordinates].(bool); ok {
		dest.Coordinates = v
	}
	if v, ok := patch[utils.KeyPopulation].(bool); ok {
		dest.Population = v
	}
	if v, ok := patch[utils.KeyArea].(bool); ok {
		dest.Area = v
	}
	if v, ok := patch[utils.KeyTargetCurrencies].([]interface{}); ok {
		var currencies []string
		for _, item := range v {
			if currency, ok := item.(string); ok {
				currencies = append(currencies, currency)
			}
		}
		dest.TargetCurrencies = currencies
	}

	log.Println("applyFeaturePatch - updated config:", dest)
}

// DeleteRegistrationByID removes a dashboard config by ID and triggers a DELETE webhook event.
func DeleteRegistrationByID(ctx context.Context, id string) error {
	config, err := db.GetDashboardConfigByID(ctx, id)
	if err != nil {
		return fmt.Errorf(utils.ErrConfigNotFoundByID, err)
	}

	if err := db.DeleteDashboardConfig(ctx, id); err != nil {
		return fmt.Errorf(utils.ErrFirestoreDeleteFailed, err)
	}

	TriggerWebhooks(utils.EventDelete, config.ISOCode)
	return nil
}
