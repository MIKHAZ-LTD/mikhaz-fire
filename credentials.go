package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2/google"
)

// Account represents a parsed Google or Firebase authenticated account.
type Account struct {
	Email        string    `json:"email"`
	Source       string    `json:"source"` // "firebase", "gcloud", "service_account"
	RefreshToken string    `json:"-"`
	AccessToken  string    `json:"-"`
	ExpiresAt    time.Time `json:"-"`
	ClientID     string    `json:"-"`
	ClientSecret string    `json:"-"`
	IsServiceAcc bool      `json:"is_service_account"`
	ServiceJSON  string    `json:"-"` // raw service account JSON
}

type CredentialManager struct {
	mu       sync.Mutex
	accounts map[string]*Account // keyed by Email
}

func NewCredentialManager() *CredentialManager {
	return &CredentialManager{
		accounts: make(map[string]*Account),
	}
}

// RefreshTokenResponse is the payload returned by Google OAuth2 token endpoint.
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// GetAccessToken returns a valid access token for the given email, refreshing it if necessary.
func (cm *CredentialManager) GetAccessToken(email string) (string, error) {
	cm.mu.Lock()
	acc, ok := cm.accounts[email]
	cm.mu.Unlock()

	if !ok {
		return "", fmt.Errorf("account not found: %s", email)
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	// If token is still valid (with a 2-minute buffer), return it
	if acc.AccessToken != "" && time.Now().Before(acc.ExpiresAt.Add(-2*time.Minute)) {
		return acc.AccessToken, nil
	}

	if acc.IsServiceAcc {
		token, err := cm.refreshServiceAccountToken(acc)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	// For standard OAuth user accounts
	token, err := cm.refreshUserToken(acc)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (cm *CredentialManager) refreshUserToken(acc *Account) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", acc.RefreshToken)
	data.Set("client_id", acc.ClientID)
	if acc.ClientSecret != "" {
		data.Set("client_secret", acc.ClientSecret)
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", fmt.Errorf("failed to make refresh request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var res RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to decode refresh response: %v", err)
	}

	acc.AccessToken = res.AccessToken
	acc.ExpiresAt = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)
	return acc.AccessToken, nil
}

// ServiceAccountKey represents Google Service Account JSON structure.
type ServiceAccountKey struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

func (cm *CredentialManager) refreshServiceAccountToken(acc *Account) (string, error) {
	return cm.refreshServiceAccountTokenOAuth2(acc)
}

func (cm *CredentialManager) refreshServiceAccountTokenOAuth2(acc *Account) (string, error) {
	config, err := google.JWTConfigFromJSON([]byte(acc.ServiceJSON), "https://www.googleapis.com/auth/cloud-platform", "https://www.googleapis.com/auth/firebase")
	if err != nil {
		return "", err
	}
	ts := config.TokenSource(context.Background())
	tok, err := ts.Token()
	if err != nil {
		return "", err
	}
	acc.AccessToken = tok.AccessToken
	acc.ExpiresAt = tok.Expiry
	return tok.AccessToken, nil
}

// We will implement scan local credentials.
func (cm *CredentialManager) ScanLocalCredentials() []Account {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	// 1. Scan Firebase CLI config: C:\Users\<username>\.config\configstore\firebase-tools.json
	firebaseConfigPath := filepath.Join(userHome, ".config", "configstore", "firebase-tools.json")
	if fileExists(firebaseConfigPath) {
		cm.parseFirebaseConfig(firebaseConfigPath)
	}

	// 2. Scan gcloud legacy credentials: C:\Users\<username>\AppData\Roaming\gcloud\legacy_credentials
	appData := os.Getenv("APPDATA")
	if appData == "" {
		appData = filepath.Join(userHome, "AppData", "Roaming")
	}
	gcloudPath := filepath.Join(appData, "gcloud")
	legacyCredsPath := filepath.Join(gcloudPath, "legacy_credentials")
	if fileExists(legacyCredsPath) {
		cm.scanGcloudLegacyCredentials(legacyCredsPath)
	}

	// 3. Scan Application Default Credentials
	adcPath := filepath.Join(gcloudPath, "application_default_credentials.json")
	if fileExists(adcPath) {
		cm.parseAdcConfig(adcPath, "gcloud-adc")
	}

	// Return list of discovered accounts
	var list []Account
	for _, acc := range cm.accounts {
		list = append(list, *acc)
	}
	return list
}

func (cm *CredentialManager) parseFirebaseConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var config struct {
		User struct {
			Email string `json:"email"`
		} `json:"user"`
		Tokens struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresAt    int64  `json:"expires_at"` // milliseconds
		} `json:"tokens"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return
	}

	if config.User.Email != "" && config.Tokens.RefreshToken != "" {
		acc := &Account{
			Email:        config.User.Email,
			Source:       "firebase",
			RefreshToken: config.Tokens.RefreshToken,
			AccessToken:  config.Tokens.AccessToken,
			ExpiresAt:    time.UnixMilli(config.Tokens.ExpiresAt),
			ClientID:     "563584335869-fgrhgmd47bqnekij5i8b5pr03ho849e6.apps.googleusercontent.com", // Firebase CLI Client ID
		}
		cm.accounts[acc.Email] = acc
	}
}

func (cm *CredentialManager) scanGcloudLegacyCredentials(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		email := entry.Name()
		adcPath := filepath.Join(dir, email, "adc.json")
		if fileExists(adcPath) {
			cm.parseAdcConfig(adcPath, email)
		}
	}
}

func (cm *CredentialManager) parseAdcConfig(path string, fallbackEmail string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var adc struct {
		Type         string `json:"type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RefreshToken string `json:"refresh_token"`
		ClientEmail  string `json:"client_email"`
	}

	if err := json.Unmarshal(data, &adc); err != nil {
		return
	}

	email := fallbackEmail
	if adc.ClientEmail != "" {
		email = adc.ClientEmail
	}

	if adc.Type == "service_account" {
		acc := &Account{
			Email:        email,
			Source:       "service_account",
			IsServiceAcc: true,
			ServiceJSON:  string(data),
		}
		cm.accounts[acc.Email] = acc
	} else if adc.RefreshToken != "" {
		emailClean := email
		if !strings.Contains(emailClean, "@") {
			emailClean = "gcloud-user"
		}
		acc := &Account{
			Email:        emailClean,
			Source:       "gcloud",
			RefreshToken: adc.RefreshToken,
			ClientID:     adc.ClientID,
			ClientSecret: adc.ClientSecret,
		}
		cm.accounts[acc.Email] = acc
	}
}

func (cm *CredentialManager) AddServiceAccount(jsonContent string) (string, error) {
	var key ServiceAccountKey
	if err := json.Unmarshal([]byte(jsonContent), &key); err != nil {
		return "", fmt.Errorf("invalid service account JSON: %v", err)
	}

	if key.ClientEmail == "" || key.PrivateKey == "" {
		return "", fmt.Errorf("missing client_email or private_key in service account JSON")
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	acc := &Account{
		Email:        key.ClientEmail,
		Source:       "service_account",
		IsServiceAcc: true,
		ServiceJSON:  jsonContent,
	}
	cm.accounts[acc.Email] = acc
	return acc.Email, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

