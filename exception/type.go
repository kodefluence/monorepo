package exception

// Type of exception
type Type uint

const (
	// Unexpected throwed when there is unexpected error that not specified yet here in exception type
	Unexpected Type = iota
	// NotFound throwed when there is unexpected data missing
	NotFound
	// Duplicated throwed when there is unexpected duplicated data
	Duplicated
	// BadInput throwed when there is unexpected input received from the caller
	BadInput
)

func (t Type) String() string {
	return []string{
		"unexpected",
		"not found",
		"duplicated",
		"bad input",
	}[t]
}
