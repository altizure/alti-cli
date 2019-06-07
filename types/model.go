package types

// Model represents the gql 'ImportedModel' type.
type Model struct {
	ID       string
	State    string
	Name     string
	Filename string
	Error    []string
}
