package utils

import (
	"encoding/json"
	"fmt"
	"github.com/enescakir/emoji"
	"os"
)

// Settings Global variable to hold the loaded configuration
var Settings MainCLIConfig

func init() {
	Settings = defaultConfig()
	err := LoadConfig(Settings.General.MainConfigFilePath)
	if err != nil {
		FatalPrint("Unable to load main config file at path %s %v", Settings.General.MainConfigFilePath, emoji.FileCabinet)
	}
}

// GitHubAuthConfig holds configuration for OAuth commands
type GitHubAuthConfig struct {
	GitHubClientID     string `json:"github_client_id"`
	GitHubClientSecret string `json:"github_client_secret"`
	RedirectURL        string `json:"redirect_url"`
	ServerPort         int    `json:"server_port"`
	AuthFilePath       string `json:"auth_file_path"`
	CallbackTimeout    int    `json:"callback_timeout"`
}

// GeneralConfig holds general configuration for the CLI application
type GeneralConfig struct {
	MainConfigFilePath string `json:"main_config_file_path"`
}

// defaultConfig returns a MainCLIConfig struct with default values
func defaultConfig() MainCLIConfig {
	return MainCLIConfig{
		General: GeneralConfig{
			MainConfigFilePath: "./config.json",
		},
		GitHubAuthConfig: GitHubAuthConfig{
			RedirectURL:     "http://localhost:%d/callback",
			ServerPort:      8080,
			AuthFilePath:    ".auth.json",
			CallbackTimeout: 20,
		},
	}
}

// MainCLIConfig holds all configurations for the CLI application
type MainCLIConfig struct {
	General          GeneralConfig    `json:"general"`
	GitHubAuthConfig GitHubAuthConfig `json:"github_auth"`
}

// LoadConfig loads the configuration from the specified file
func LoadConfig(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Settings)
	if err != nil {
		return fmt.Errorf("could not decode config JSON: %v", err)
	}

	return nil
}
