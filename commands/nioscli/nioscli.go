package nioscli

import (
	"os"

	"github.com/dev-vinicius-andrade/nioscli/commands/nioscli/create"

	"github.com/dev-vinicius-andrade/nioscli/commands/nioscli/templates"
	"github.com/dev-vinicius-andrade/nioscli/helpers"
	"github.com/dev-vinicius-andrade/nioscli/types"
	"github.com/dev-vinicius-andrade/nioscli/types/context"
	"github.com/spf13/cobra"
)

type commandDefinition struct {
	cmd *cobra.Command
	//information *context.DotfilesInformationContext
	context *context.NiOsContext
}

func (c *commandDefinition) createCompletitionCommand() *cobra.Command {
	command := &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Generate completion script",
		Long:                  "To load completions",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(c.cmd.Args),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
	return command
}
func (c *commandDefinition) createCommandDefinition() *cobra.Command {

	command := &cobra.Command{
		Use:              "nioscli",
		Short:            "nioscli is a CLI tool help bootstrap NixOs.",
		Long:             `The goal of this tool is to help you to bootstrap your NixOs using a simple CLI tool`,
		TraverseChildren: true,
		Run:              c.runCommand,
	}

	c.context = createContext()
	c.cmd = command
	c.defineFlags()
	c.defineSubCommands()
	return command
}

func (c *commandDefinition) defineSubCommands() {
	create.CreateCommand(c.context, c.cmd)
	templates.CreateCommand(c.context, c.cmd)
	c.cmd.AddCommand(c.createCompletitionCommand())
}

func (c *commandDefinition) runCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		helpers.CobraHelper.ShowHelp(cmd, &types.CobraHelpOptions{Title: "Nixos installer command command"})
		return
	}
}

func (c *commandDefinition) defineFlags() {
}

func CreateCommand() *cobra.Command {
	c := commandDefinition{}
	return c.createCommandDefinition()
}
