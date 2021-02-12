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
	// Unauthorized throwed when there is unauthorized access from the caller
	Unauthorized
	// Forbidden throwd when there is unexpected access from the caller
	Forbidden
)

func (t Type) String() string {
	return []string{
		"unexpected",
		"not found",
		"duplicated",
		"bad input",
		"unauthorized",
		"forbidden",
	}[t]
}
