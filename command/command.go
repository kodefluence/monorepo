package command

import (
	"github.com/spf13/cobra"
)

//go:generate mockgen -source=./command.go -destination=./mock/mock.go -package mock

// Scaffold use for standard of creating CLI command
type Scaffold interface {
	Use() string
	Example() string
	Short() string
	Run(args []string)
}

// Command manage all command in monorepo
type Command struct {
	rootCmd *cobra.Command
}

// Fabricate monorepo command, it just get the first parameters of configs.
func Fabricate(configs ...Config) *Command {
	config := Config{
		Name:  "monorepo",
		Short: "A placeholder you need to actually change it",
	}

	if len(configs) > 0 {
		config.Name = configs[0].Name
		config.Short = configs[0].Short
	}

	return &Command{
		rootCmd: &cobra.Command{
			Use:     config.Name,
			Short:   config.Short,
			Example: config.Name,
			Run: func(cmd *cobra.Command, args []string) {
				_ = cmd.Help()
			},
		},
	}
}

// SetArgs set argument for command line interface
func (c *Command) SetArgs(args []string) {
	c.rootCmd.SetArgs(args)
}

// Execute command line interface
func (c *Command) Execute() error {
	return c.rootCmd.Execute()
}

// InjectCommand inject new command into command list
func (c *Command) InjectCommand(scaffolds ...Scaffold) {
	for _, scaffold := range scaffolds {
		// Intendedly assign this variable
		scaffoldRunFunction := scaffold.Run

		cmd := &cobra.Command{
			Use:     scaffold.Use(),
			Short:   scaffold.Short(),
			Example: scaffold.Example(),
			Run: func(cmd *cobra.Command, args []string) {
				scaffoldRunFunction(args)
			},
		}
		c.rootCmd.AddCommand(cmd)
	}
}
