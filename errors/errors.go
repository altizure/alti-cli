package errors

const (
	// ErrNoConfig is returned when config file is not found.
	ErrNoConfig AppError = "app: no config"
	// ErrNotLogin is returned when user is not login.
	ErrNotLogin AppError = "app: not login"
	// ErrProfileNotFound is returned when the queried profile is not found.
	ErrProfileNotFound ConfigError = "config: profile not found"
	// ErrOffline is returned when the server is offline.
	ErrOffline ServerError = "server: offline"
)

// AppError is the application specific error.
type AppError string

func (e AppError) Error() string {
	return string(e)
}

// ConfigError is the config specific error.
type ConfigError string

func (e ConfigError) Error() string {
	return string(e)
}

// ServerError is the server specific error.
type ServerError string

func (e ServerError) Error() string {
	return string(e)
}
