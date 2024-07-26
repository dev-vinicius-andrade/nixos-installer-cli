package configuration

import (
	"github.com/dev-vinicius-andrade/nioscli/commands/nioscli/create/configuration/hosts"
	"github.com/dev-vinicius-andrade/nioscli/helpers"
	"github.com/dev-vinicius-andrade/nioscli/types/context"
	"github.com/spf13/cobra"
)

type commandDefinition struct {
	cmd           *cobra.Command
	context       *context.NiOsContext
	parentCommand *cobra.Command
	//Add more command definition variables here

}

func (c *commandDefinition) runCommand(cmd *cobra.Command, args []string) {
}
func (c *commandDefinition) createCommandDefinition() *cobra.Command {
	command := &cobra.Command{
		Use:   "configurations",
		Short: "Configuration utility",
		Long:  "Configuration utility",
		Run: func(cmd *cobra.Command, args []string) {
			c.runCommand(cmd, args)
		},
	}
	c.cmd = command
	c.defineFlags()
	c.defineSubCommands()
	helpers.CobraHelper.AddCommandToParent(command, c.parentCommand)
	return command
}
func (c *commandDefinition) defineSubCommands() {
	hosts.CreateCommand(c.context, c.cmd)
	// Add subcommands here
}
func (c *commandDefinition) defineFlags() {
	// Add flags here
}
func CreateCommand(context *context.NiOsContext, parentCommand *cobra.Command) *cobra.Command {
	command := &commandDefinition{
		context:       context,
		parentCommand: parentCommand,
		//Add your command definition extra default values here
	}
	return command.createCommandDefinition()
}
