package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type Suggestion struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AnalyseResponse struct {
	AnalysisID  string       `json:"analysisId"`
	PanelURL    string       `json:"panelUrl"`
	Suggestions []Suggestion `json:"suggestions"`
}

func loadConfig() (apiKey string, apiURL string, err error) {
	configDir, _ := os.UserConfigDir()

	keyData, err := os.ReadFile(configDir + "/repomind/config")
	if err != nil {
		return "", "", fmt.Errorf("API key not found. Run 'repomind init <api-key>' first")
	}

	apiURL = os.Getenv("REPOMIND_API_URL")
	if apiURL == "" {
		urlData, err := os.ReadFile(configDir + "/repomind/url")
		if err != nil {
			apiURL = "https://repomind-api-htom.onrender.com"
		} else {
			apiURL = string(urlData)
		}
	}

	return string(keyData), apiURL, nil
}

func Analyse(files []map[string]string, context string, projectName string) (*AnalyseResponse, error) {
	apiKey, apiURL, err := loadConfig()
	if err != nil {
		return nil, err
	}

	client := resty.New()

	body := map[string]any{
		"files":       files,
		"context":     context,
		"projectName": projectName,
	}

	resp, err := client.R().
		SetHeader("x-api-key", apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(apiURL + "/analyse")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 429 {
		var errBody map[string]any
		json.Unmarshal(resp.Body(), &errBody)
		return nil, fmt.Errorf("rate_limit: %v", errBody["error"])
	}

	var result AnalyseResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
