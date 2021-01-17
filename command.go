package monorepo

//go:generate mockgen -source=./command.go -destination=./monomock/command_mock.go -package monomock

// Command CLI command bearer of monorepo
type Command interface {
	InjectCommand(scaffold ...CommandScaffold)
}

// CommandScaffold use for standard of creating CLI command
type CommandScaffold interface {
	Use() string
	Example() string
	Short() string
	Run(args []string)
}
