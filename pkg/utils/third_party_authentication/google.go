package third_party_authentication

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}

func GetGoogleIDToken(
	authorizationCode string,
	platform consts.Platform,
) (idToken string, err error) {
	// 1. Determine client_id based on platform
	var clientId string
	switch platform {
	case consts.WEB:
		clientId = global.Config.GoogleSetting.WebClientId
	case consts.ANDROID:
		clientId = global.Config.GoogleSetting.AndroidClientId
	case consts.IOS:
		clientId = global.Config.GoogleSetting.IosClientId
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}

	// 2. Prepare payload
	payload := url.Values{}
	payload.Set("code", authorizationCode)
	payload.Set("client_id", clientId)
	payload.Set("client_secret", global.Config.GoogleSetting.SecretId)
	payload.Set("redirect_uri", global.Config.GoogleSetting.RedirectUrl)
	payload.Set("grant_type", "authorization_code")

	// 3. Create HTTP request
	req, err := http.NewRequest("POST", global.Config.GoogleSetting.GoogleTokensUrl, strings.NewReader(payload.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 4. Send the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Google: %w", err)
	}
	defer resp.Body.Close()

	// 5. Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google API returned error: %d %s", resp.StatusCode, resp.Status)
	}

	// 6. Decode response
	var tokenResponse GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("failed to parse response from Google: %w", err)
	}

	// 7. Validate ID token
	if tokenResponse.IDToken == "" {
		return "", errors.New("google response does not contain an ID token")
	}

	return tokenResponse.IDToken, nil
}
