package db

// Config carry all database config in a single struct
type Config struct {
	Username string
	Password string
	Host     string

	// If port is empty it will return 3306 instead
	Port string

	Name string
}
