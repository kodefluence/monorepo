package exception

// A Config for Exception creation
type Config struct {
	title         string
	detail        string
	exceptionType Type
}

// Option when fabricating Exception
type Option func(*Config)

// WithTitle fabricate exception with title
func WithTitle(title string) Option {
	return func(c *Config) {
		c.title = title
	}
}

// WithDetail fabricate exception with detail
func WithDetail(detail string) Option {
	return func(c *Config) {
		c.detail = detail
	}
}

// WithType fabricate exception with exception type
func WithType(exceptionType Type) Option {
	return func(c *Config) {
		c.exceptionType = exceptionType
	}
}
