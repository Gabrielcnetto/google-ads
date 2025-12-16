package gads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	gads "netto/controllers/gads/credentials"
	"netto/models"
)

func GetAcessToken() (token string, err error) {
	body := models.Body{
		ClientId:     gads.C_I,
		ClientSecret: gads.C_S,
		RefreshToken: gads.R_T,
		GrantType:    "refresh_token",
	}
	bodyJson, _ := json.Marshal(body)
	url := "https://oauth2.googleapis.com/token"
	r, err := http.Post(url, "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	responseBody, err := io.ReadAll(r.Body)

	if err != nil {
		return "", err
	}
	var decodedData map[string]any
	if err := json.Unmarshal(responseBody, &decodedData); err != nil {
		return "", err
	}
	accessToken, ok := decodedData["access_token"].(string)
	if !ok || accessToken == "" {
		return "", fmt.Errorf("access_token não encontrado ou inválido")
	}
	return accessToken, nil
}
