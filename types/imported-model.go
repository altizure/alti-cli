package types

// ImportedModel represents the gql 'ImportedModel' type.
type ImportedModel struct {
	ID       string
	State    string
	Name     string
	Filename string
	Error    []string
}
