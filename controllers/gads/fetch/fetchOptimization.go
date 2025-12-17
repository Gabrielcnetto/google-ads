package gads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	gads "netto/controllers/gads/credentials"
	"strings"
)

func FetchOptimization(accessToken string, customerID string) (body map[string]interface{}, err error) {
	url := fmt.Sprintf("https://googleads.googleapis.com/v22/customers/%v/googleAds:search", strings.Replace(customerID, "-", "", -1))
	requestBody := []byte(`
	{
		"query": "SELECT customer.id,customer.descriptive_name,customer.optimization_score FROM customer"
	}
	`)
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
	var decodedData map[string]interface{}
	if err := json.Unmarshal(responseBody, &decodedData); err != nil {
		return nil, err
	}

	results, ok := decodedData["results"]
	if !ok {
		return nil, fmt.Errorf("campo 'results' não encontrado")
	}
	resultsSlice, ok := results.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'results' não é uma lista")
	}
	firstItem, ok := resultsSlice[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("item 0 não é um objeto JSON")
	}
	customer, ok := firstItem["customer"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Erro ao pegar o user account")
	}
	optValue, ok := customer["optimizationScore"].(float64)
	if !ok {
		optValue = 0
	}

	return map[string]interface{}{
		"optimazion_score": optValue,
	}, nil
}
