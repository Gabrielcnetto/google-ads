package gads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	gads "netto/controllers/gads/credentials"
	"strconv"
	"strings"
)

func FetchAccountDatas(accessToken string, customerID string, month int, day int, year int) (body map[string]interface{}, err error) {

	url := fmt.Sprintf("https://googleads.googleapis.com/v22/customers/%v/googleAds:search", strings.Replace(customerID, "-", "", -1))
	query := fmt.Sprintf("SELECT campaign.name, metrics.all_conversions, metrics.cost_micros, metrics.impressions, metrics.clicks FROM campaign WHERE campaign.status IN ('ENABLED') AND segments.date BETWEEN '%d-%02d-01' AND '%d-%02d-%02d'", year, month, year, month, day)
	requestBody := []byte(fmt.Sprintf(`
	{
		"query": "%s"
	}
	`, query))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("developer-token", gads.D_T)
	req.Header.Set("login-customer-id", gads.MCC_ID)
	bearer := fmt.Sprintf("Bearer %v", accessToken)
	req.Header.Set("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var decodedData map[string]any
	if err := json.Unmarshal(responseBody, &decodedData); err != nil {
		return nil, err
	}
	results, ok := decodedData["results"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Erro ao pegar apenas result de métricas")
	}

	var CampaignInfs []map[string]interface{}
	var totalCost = 0.0
	for _, item := range results {
		castItem, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("Erro: item não é do tipo map[string]interface{}")
			continue
		}

		campaign, ok := castItem["campaign"].(map[string]interface{})
		if !ok {
			fmt.Println("Erro ao acessar 'campaign'")
			continue
		}

		campaignName, _ := campaign["name"].(string)

		metrics, ok := castItem["metrics"].(map[string]interface{})
		if !ok {
			fmt.Println("Erro ao acessar 'metrics'")
			continue
		}

		clicks, ok := metrics["clicks"].(string)
		if !ok {
			clicks = "0"
		}

		conversions, ok := metrics["allConversions"].(float64)
		if !ok {
			conversions = 0
		}

		costMicros, ok := metrics["costMicros"].(string)
		if !ok {
			costMicros = "0"
		}

		floatNum, err := strconv.ParseFloat(costMicros, 64)
		if err != nil {
			floatNum = 0
		}

		convertedClicks, err := strconv.Atoi(clicks)
		if err != nil {
			convertedClicks = 0
		}

		impressions, ok := metrics["impressions"].(string)
		if !ok {
			impressions = "0"
		}

		convertedImpr, err := strconv.Atoi(impressions)
		if err != nil {
			convertedImpr = 0
		}

		item := map[string]interface{}{
			"clicks":       convertedClicks,
			"cost":         floatNum / 1000000,
			"impressions":  convertedImpr,
			"campaignName": campaignName,
			"conversions":  conversions,
		}
		totalCost += (floatNum / 1000000)
		CampaignInfs = append(CampaignInfs, item)
	}

	return map[string]interface{}{
		"campanhas": CampaignInfs,
		"cost":      totalCost,
	}, nil

}
