package create

import (
	"github.com/dev-vinicius-andrade/nioscli/commands/nioscli/create/configuration"
	"github.com/dev-vinicius-andrade/nioscli/commands/nioscli/create/disko"
	"github.com/dev-vinicius-andrade/nioscli/commands/nioscli/create/hardware"
	"github.com/dev-vinicius-andrade/nioscli/helpers"
	"github.com/dev-vinicius-andrade/nioscli/types"
	"github.com/dev-vinicius-andrade/nioscli/types/context"
	"github.com/spf13/cobra"
)

type commandDefinition struct {
	cmd           *cobra.Command
	context       *context.NiOsContext
	parentCommand *cobra.Command
	Version       bool
}

func (c *commandDefinition) runCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		helpers.CobraHelper.ShowHelp(cmd, &types.CobraHelpOptions{Title: "Nixos installer command"})
		return
	}
}
func (c *commandDefinition) printVersion() {
	//fmt.Printf("Dotfiles version: %s\n", c.context.DotfilesInformation.Version)

}
func (c *commandDefinition) createCommandDefinition() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "Creates files/folders or disk partitions",
		Long:  `Create files/folders or disk partitions, checkout the flags to see the options`,
		Run: func(cmd *cobra.Command, args []string) {
			c.runCommand(cmd, args)
		},
	}
	c.cmd = command

	c.defineFlags()
	c.defineSubCommands()
	helpers.CobraHelper.AddCommandToParent(c.cmd, c.parentCommand)
	return &c.context.NiOsCmd
}
func (c *commandDefinition) defineFlags() {
	//c.cmd.Flags().BoolVarP(&c.Version, "version", "v", false, "Get the current version of the dotfiles tool")
}
func (c *commandDefinition) defineSubCommands() {
	disko.CreateCommand(c.context, c.cmd)
	hardware.CreateCommand(c.context, c.cmd)
	configuration.CreateCommand(c.context, c.cmd)
}
func CreateCommand(context *context.NiOsContext, parentCommand *cobra.Command) *cobra.Command {
	c := commandDefinition{
		context:       context,
		parentCommand: parentCommand,
	}
	return c.createCommandDefinition()
}
