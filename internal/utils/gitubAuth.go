package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var githubUser GitHubUser
var oauth2Config *oauth2.Config
var httpClient *http.Client

type GitHubUser struct {
	Username string `json:"user_name" validate:"required"`
	Token    string `json:"access_token" validate:"required"`
	ClientId string `json:"client_id" validate:"required"`
}

// User represents the GitHub user data
type User struct {
	Login string `json:"login"`
}

func init() {
	httpClient = &http.Client{}
	githubUser = GitHubUser{}
}

func initOAuthConfig() {
	oauth2Config = &oauth2.Config{
		ClientID:     Settings.GitHubAuthConfig.GitHubClientID,
		ClientSecret: Settings.GitHubAuthConfig.GitHubClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  fmt.Sprintf(Settings.GitHubAuthConfig.RedirectURL, Settings.GitHubAuthConfig.ServerPort),
		Scopes:       []string{"repo", "user"},
	}
}

// openBrowser Opens the default configured browser for the current user.
func openBrowser(url string) {
	err := exec.Command("open", url).Start()
	if err != nil {
		FatalPrint("Error opening browser: %v\n", err)
	}
}

func saveAuthInfo() error {
	file, err := os.Create(Settings.GitHubAuthConfig.AuthFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(githubUser)
	if err != nil {
		return err
	}
	return nil
}

func removeAuthInfo() error {
	err := os.Remove(Settings.GitHubAuthConfig.AuthFilePath)
	if err != nil {
		return err
	}
	return nil
}

// loadAuthToken loads the github auth Token from the local file.
func loadAuthToken() error {
	// Read the JSON file
	data, err := os.ReadFile(Settings.GitHubAuthConfig.AuthFilePath)
	if err != nil {
		return &CustomError{Code: AuthFileNotExistsError, Message: "file not found"}
	}

	// Unmarshal the JSON into a Config struct
	if err := json.Unmarshal(data, &githubUser); err != nil {
		return &CustomError{Code: AuthFileNotValidJsonError, Message: "invalid json context in auth file."}
	}

	// Validate the Config struct
	validate := validator.New()
	if err := validate.Struct(githubUser); err != nil {
		return &CustomError{Code: AuthFileInvalidError, Message: "invalid user auth data"}
	}
	return nil
}

// setupCallbackLogic - Creates an HTTP endpoint that receives the 0Auth callback after authentication.
// Returns:
// - chan int: Go routine channel that signals the authentication status.
func setupCallbackLogic(authState string) chan int {
	callbackChannel := make(chan int)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != authState {
			http.Error(w, "State did not match", http.StatusBadRequest)
			FatalPrint("Login request state did not match the request's state param.")
		}
		code := r.URL.Query().Get("code")
		token, err := oauth2Config.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, "Failed to authenticate current user", http.StatusForbidden)
			FatalPrint("Failed to authenticate current user")
		}
		githubUser.Token = token.AccessToken
		username, err := getGitHubUsername()
		if err != nil {
			http.Error(w, "Failed to fetch GitHub user info: "+err.Error(), http.StatusInternalServerError)
			FatalPrint("Failed to fetch GitHub user info")
		}
		githubUser.Username = username
		githubUser.ClientId = Settings.GitHubAuthConfig.GitHubClientID
		err = saveAuthInfo()
		// Send ascii art payload as web response
		_, err = w.Write([]byte(MainArt + SuccessAscii))
		if err != nil {
			http.Error(w, "Failed to authenticate user: "+err.Error(), http.StatusInternalServerError)
		}
		// Signal that the authentication is complete.
		close(callbackChannel)
	})
	return callbackChannel
}

func getGitHubUsername() (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+githubUser.Token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}
	return user.Login, nil
}

// generateStateParam generates a 16-byte random string to send as the CSRF protection param.
// Returns:
// - string: Url encoded random string
func generateStateParam() (string, error) {
	// OAuth2 state param generate.
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// PerformLogin initiates the OAuth 2.0 login process by opening the user's browser to GitHub's
// authorization URL and starting a local HTTP server to handle the OAuth callback.
// Returns:
// - error: Returns an error if any step in the login process fails.
func PerformLogin() error {
	initOAuthConfig()
	state, err := generateStateParam()
	if err != nil {
		return fmt.Errorf("failed to generate state: %w", err)
	}
	authURL := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	openBrowser(authURL)
	callbackChannel := setupCallbackLogic(state)
	InfoPrint("Please login to your Github account using the opened browser window")

	// Start the server in a goroutine
	server := &http.Server{Addr: fmt.Sprintf(":%d", Settings.GitHubAuthConfig.ServerPort)}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
			close(callbackChannel) // Signal completion on error as well
		}
	}()
	// Create a timeout channel
	timeout := time.After(time.Duration(Settings.GitHubAuthConfig.CallbackTimeout) * time.Second)

	// Wait for either the callback or the timeout
	select {
	case <-callbackChannel:
		SuccessPrint("Welcome, %s you have been authenticated!\n", githubUser.Username)
	case <-timeout:
		server.Close() // Stop the server
		ErrorPrint("Timeout waiting for Github authentication, exiting.. %v", emoji.SadButRelievedFace)
		os.Exit(1)
	}
	return nil
}

func Logout() {
	// Check if the user is even logged in
	if !CheckLoginStatus() {
		WarningPrint("User isn't logged in, no reason to log out %v", emoji.Teapot)
		return
	}
	url := fmt.Sprintf("https://api.github.com/applications/%s/grant", Settings.GitHubAuthConfig.GitHubClientID)
	// Create the request payload
	payload := map[string]string{
		"access_token": githubUser.Token,
	}
	jsonPayload, err := json.Marshal(payload)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		FatalPrint("Failed to create logout request: %v\n", err)
	}

	// Set the required headers
	req.Header.Set("Authorization", "Basic "+BasicAuth(Settings.GitHubAuthConfig.GitHubClientID, Settings.GitHubAuthConfig.GitHubClientSecret))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		if resp != nil {
			FatalPrint("Failed to logout, received bad response code %d", resp.StatusCode)
		} else {
			FatalPrint("Failed to logout")
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		FatalPrint("Failed to logout, response status code: %v", resp.StatusCode)
	}
	err = removeAuthInfo()
	if err != nil {
		ErrorPrint("Failed to remove authentication file: %s", err.Error())
	}
	SuccessPrint("Successfully logged out %v", emoji.Hamburger)
}

// CheckLoginStatus checks if the current authentication Token is valid.
// Returns:
// - bool: true if the authentication Token is valid, false otherwise.
func CheckLoginStatus() bool {
	err := loadAuthToken()
	if err != nil {
		return false
	}
	// This means that a valid auth token has been loaded.
	if githubUser.Token == "" {
		return false
	}
	// Validate Token with GitHub
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		FatalPrint("Error validating current authentication Token, error: %v\n", err)
		return false
	}
	req.Header.Set("Authorization", "Bearer "+githubUser.Token)

	resp, err := httpClient.Do(req)
	if err != nil {
		FatalPrint("Error validating current authentication Token, error: %v\n", err)
	}
	defer resp.Body.Close()

	// Check if the Token is still valid
	if resp.StatusCode != http.StatusOK {
		ErrorPrint("Failed to validate current user's authentication status: %v\n", resp.StatusCode)
		return false
	}
	return true
}
