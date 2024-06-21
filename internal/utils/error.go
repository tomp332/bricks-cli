package utils

// CustomError represents an error with a code
type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

const (
	AuthFileInvalidError      = 1
	AuthFileNotExistsError    = 2
	AuthFileNotValidJsonError = 3
)
