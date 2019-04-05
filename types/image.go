package types

// Image represents the gql 'ProjectImage' type.
type Image struct {
	ID       string
	State    string
	Name     string
	Filename string
	Error    []string
}
