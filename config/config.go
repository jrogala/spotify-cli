package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigDir())

	viper.SetEnvPrefix("SPOTIFY")
	viper.AutomaticEnv()

	viper.SetDefault("redirect_uri", "http://127.0.0.1:8888/callback")

	_ = viper.ReadInConfig()
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "spotify-cli")
}

func ClientID() string {
	return viper.GetString("client_id")
}

func RedirectURI() string {
	return viper.GetString("redirect_uri")
}

func DefaultDevice() string {
	return viper.GetString("default_device")
}

func Save() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return viper.WriteConfigAs(filepath.Join(dir, "config.yaml"))
}

// Token represents stored OAuth2 tokens.
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scope        string    `json:"scope"`
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt.Add(-60 * time.Second))
}

func tokenPath() string {
	return filepath.Join(ConfigDir(), "token.json")
}

func LoadToken() (*Token, error) {
	data, err := os.ReadFile(tokenPath())
	if err != nil {
		return nil, fmt.Errorf("not authenticated, run 'spotify auth login'")
	}
	var t Token
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("corrupt token file: %w", err)
	}
	return &t, nil
}

func SaveToken(t *Token) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tokenPath(), data, 0600)
}

func DeleteToken() error {
	return os.Remove(tokenPath())
}
