package exception

// Exception wrap error into more informative structs
type Exception interface {
	Error() string
	Type() Type
	Detail() string
	Title() string
}

// Error wrap default error with exception format
type Error struct {
	config Config
	err    error
}

// Throw new exception
func Throw(err error, opts ...Option) Exception {
	var config Config

	// Default value
	config.exceptionType = Unexpected

	for _, opt := range opts {
		opt(&config)
	}

	return &Error{
		config: config,
		err:    err,
	}
}

func (e *Error) Error() string {
	return e.err.Error()
}

// Type return exception type
func (e *Error) Type() Type {
	return e.config.exceptionType
}

// Detail return exception detail
func (e *Error) Detail() string {
	return e.config.detail
}

// Title return exception title
func (e *Error) Title() string {
	return e.config.title
}
