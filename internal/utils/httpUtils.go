package utils

import "encoding/base64"

// BasicAuth function creates a basic auth header value
// Returns:
// - string: The base64 encoded Bearer token
func BasicAuth(user, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(user + ":" + password))
}
