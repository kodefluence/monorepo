package db

import "time"

// Config carry all database config in a single struct
type Config struct {
	Username string
	Password string
	Host     string

	// If port is empty it will return 3306 instead
	Port string

	Name string

	maxIdleConn     int
	maxOpenConn     int
	connMaxLifetime time.Duration
}

// Option when fabricating connection
type Option func(*Config)

// WithMaxIdleConn fabricate connection with max idle connection
func WithMaxIdleConn(maxIdleConn int) Option {
	return func(c *Config) {
		c.maxIdleConn = maxIdleConn
	}
}

// WithMaxOpenConn fabricate connection with max open connection
func WithMaxOpenConn(maxOpenConn int) Option {
	return func(c *Config) {
		c.maxOpenConn = maxOpenConn
	}
}

// WithConnMaxLifetime fabricate connection with max connection lifetime
func WithConnMaxLifetime(connMaxLifetime time.Duration) Option {
	return func(c *Config) {
		c.connMaxLifetime = connMaxLifetime
	}
}
