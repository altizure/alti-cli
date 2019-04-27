package types

// STS represents the STS credentials returned from api server.
type STS struct {
	ID       string
	Secret   string
	Token    string
	Bucket   string
	Endpoint string
	Expire   string
}
