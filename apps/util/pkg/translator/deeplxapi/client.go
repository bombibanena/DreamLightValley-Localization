package deeplxapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ddv_loc/pkg/translator/deeplxapi/config"
)

type (
	deeplxApiClient struct {
		cfg config.Config
	}

	translationRequest struct {
		Text       string `json:"text"`
		SourceLang string `json:"source_lang"`
		TargetLang string `json:"target_lang"`
	}

	translationResponse struct {
		Alternatives []string `json:"alternatives"`
		Code         int      `json:"code"`
		Data         string   `json:"data"`
		ID           int64    `json:"id"`
		Method       string   `json:"method"`
		SourceLang   string   `json:"source_lang"`
		TargetLang   string   `json:"target_lang"`
	}
)

func NewClient(cfg config.Config) *deeplxApiClient {
	return &deeplxApiClient{
		cfg: cfg,
	}
}

func (c *deeplxApiClient) Translate(text string, sourceLang, targetLang string) (string, error) {
	reqBody := translationRequest{
		Text:       text,
		SourceLang: sourceLang,
		TargetLang: targetLang,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("Error marshaling request body:", err)
	}

	req, err := http.NewRequest("POST", c.cfg.API.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Error creating request:", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.API.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error: %s", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body:", err)
	}

	var translationResponse translationResponse
	err = json.Unmarshal(body, &translationResponse)
	if err != nil {
		return "", fmt.Errorf("Error unmarshaling response body:", err)
	}

	return translationResponse.Data, nil
}
