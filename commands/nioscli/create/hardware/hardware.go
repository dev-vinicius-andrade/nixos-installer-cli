package hardware

import (
	"fmt"
	"os"
	"os/exec"

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
	destination   string
	noFileSytem   bool
	root          string
	moveFile      bool
}

func (c *commandDefinition) generateConfig(cmd *cobra.Command, args []string) {
	// get current os
	spinner := helpers.NewSpinner()
	spinner.SetMessage("Generating hardware configuration...")
	spinner.Start()
	cmdArgsBuilder := fmt.Sprintf("--root %s", c.root)
	if c.noFileSytem {
		cmdArgsBuilder += " --no-filesystems"
	}

	commandString := fmt.Sprintf("nixos-generate-config %s", cmdArgsBuilder)
	//fmt.Printf("Running command:  %s \n", commandString)
	spinner.SetMessage(fmt.Sprintf("Running command:  %s", commandString))
	command := exec.Command("/bin/sh", "-c", commandString)
	a, err := command.CombinedOutput()
	fmt.Printf("\n%s\n", a)
	//spinner.Finish()
	if err != nil {
		colors.Default.Error.Println("Error generating hardware configuration")
		colors.Default.Error.Println(err)
		color.Unset()
		os.Exit(1)
	}
	colors.Default.Success.Println("Hardware configuration generated")
	color.Unset()
}
func (c *commandDefinition) moveHardwareConfigToDestination(cmd *cobra.Command, args []string) {
	fmt.Println("Moving hardware configuration to destination")
	if !c.moveFile {
		return
	}
	// if c.destination == "" {
	// 	colors.Default.Error.Println("Destination path is required")
	// 	color.Unset()
	// 	os.Exit(1)
	// }
	info, destinationError := os.Stat(c.destination)
	if os.IsNotExist(destinationError) {
		fmt.Println("Destination path does not exist")
		return
	}
	if !info.IsDir() {
		fmt.Println("Destination path is not a directory")
		return
	}
	source := fmt.Sprintf("%s/etc/nixos/hardware-configuration.nix", c.root)
	commandString := fmt.Sprintf("mv %s %s", source, c.destination)
	command := exec.Command("/bin/sh", "-c", commandString)

	_, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return
	}

}
func (c *commandDefinition) runCommand(cmd *cobra.Command, args []string) {
	c.generateConfig(cmd, args)
	c.moveHardwareConfigToDestination(cmd, args)
}
func (c *commandDefinition) createCommandDefinition() *cobra.Command {
	command := &cobra.Command{
		Use:   "hardware",
		Short: "Hardware utility",
		Long:  `Helps you to create the hardware configuration`,
		Run: func(cmd *cobra.Command, args []string) {
			c.runCommand(cmd, args)
		},
	}
	c.cmd = command
	c.defineFlags()
	helpers.CobraHelper.AddCommandToParent(command, c.parentCommand)
	return &c.context.NiOsCmd
}
func (c *commandDefinition) defineFlags() {
	c.cmd.Flags().StringVar(&c.destination, "destination", "", "Sets the destination path, basically it will move the /etc/nixos/hardware-configuration.nix to the destination path")
	c.cmd.Flags().BoolVar(&c.noFileSytem, "no-filesystem", false, "Flag to not create the filesystem")
	c.cmd.Flags().StringVar(&c.root, "root", "/mnt", "Sets the root path")
	c.cmd.Flags().BoolVar(&c.moveFile, "move-file", false, "Flag to move the hardware configuration file to the destination path")
	c.cmd.MarkFlagsRequiredTogether("move-file", "destination")
}
func CreateCommand(context *context.NiOsContext, parentCommand *cobra.Command) *cobra.Command {
	c := commandDefinition{
		context:       context,
		parentCommand: parentCommand,
		//path:          "./disko.nix",
	}
	return c.createCommandDefinition()
}
