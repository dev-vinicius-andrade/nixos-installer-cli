package disko

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dev-vinicius-andrade/nioscli/commands/nioscli/create/disko/config"
	"github.com/dev-vinicius-andrade/nioscli/helpers"
	"github.com/dev-vinicius-andrade/nioscli/types/colors"
	"github.com/dev-vinicius-andrade/nioscli/types/context"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type commandDefinition struct {
	cmd           *cobra.Command
	context       *context.NiOsContext
	parentCommand *cobra.Command
	path          string
	repository    string
	mode          string
	showTrace     bool
}

func (c *commandDefinition) runCommand(cmd *cobra.Command, args []string) {
	// get current os
	fmt.Printf("Running disko command with path: %s\n", c.path)
	//successColor := color.New(color.FgGreen, color.Bold)
	//errorColor := color.New(color.FgRed, color.Bold)
	info, destinationError := os.Stat(c.path)
	if os.IsNotExist(destinationError) {
		fmt.Printf("The path %s does not exist\n", c.path)
		return
	}
	if info.IsDir() {
		fmt.Printf("The path %s is a directory \n", c.path)
		return
	}
	commandString := fmt.Sprintf("nix --experimental-features \"nix-command flakes\" run %s -- --mode %s $(realpath %s)", c.repository, c.mode, c.path)
	if c.showTrace {
		commandString += " --show-trace"
	}
	spinner := helpers.NewSpinner()
	spinner.SetMessage("Running disko command")
	spinner.Start()
	spinner.SetMessage(fmt.Sprintf("Running command: %s", commandString))

	command := exec.Command("/bin/sh", "-c", commandString)

	output, err := command.CombinedOutput()
	spinner.Finish()

	if err != nil {
		colors.Default.Error.Println("Error running disko command")
		fmt.Println(err)
		color.Unset()
		os.Exit(1)
	}

	fmt.Println(string(output))
	colors.Default.Success.Println("Disko command ran successfully")
}
func (c *commandDefinition) createCommandDefinition() *cobra.Command {
	command := &cobra.Command{
		Use:   "disko",
		Short: "Disko utility",
		Long:  `Setups disko, a tool that helps you to format and partition your disks`,
		Run: func(cmd *cobra.Command, args []string) {
			c.runCommand(cmd, args)
		},
	}
	c.cmd = command
	c.defineFlags()
	c.defineSubCommands()
	helpers.CobraHelper.AddCommandToParent(command, c.parentCommand)
	return c.parentCommand
}
func (c *commandDefinition) defineSubCommands() {
	config.CreateCommand(c.context, c.cmd)
}
func (c *commandDefinition) defineFlags() {
	c.cmd.Flags().StringVar(&c.path, "path", "./disko.nix", "Get all versions tags of the dotfiles tool")
	c.cmd.Flags().StringVar(&c.repository, "repository", "github:nix-community/disko", "Disko repository url")
	c.cmd.Flags().StringVar(&c.mode, "mode", "disko", "Disko mode")
	c.cmd.Flags().BoolVar(&c.showTrace, "show-trace", false, "Show trace")

}
func CreateCommand(context *context.NiOsContext, parentCommand *cobra.Command) *cobra.Command {
	c := commandDefinition{
		context:       context,
		parentCommand: parentCommand,
		path:          "./disko.nix",
	}
	return c.createCommandDefinition()
}
