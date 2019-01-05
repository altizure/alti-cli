package gql

const (
	// ErrNoConfig is returned when config file is not found.
	ErrNoConfig AppError = "app: no config"
	// ErrNotLogin is returned when user is not login
	ErrNotLogin AppError = "app: not login"
)

// AppError is the application specific error.
type AppError string

func (e AppError) Error() string {
	return string(e)
}
