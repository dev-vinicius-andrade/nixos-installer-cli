package templates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	source        string
	destination   string
}

func (c *commandDefinition) validateFileInfo(parameterName string, fileInfo fs.FileInfo, err *error) bool {

	if os.IsNotExist(*err) {
		colors.Default.Error.Println(fmt.Sprintf("%s directory does not exist", parameterName))
		color.Unset()
		return false
	}
	if !fileInfo.IsDir() {
		colors.Default.Error.Println(fmt.Sprintf("%s  is not a directory", parameterName))
		color.Unset()
		return false
	}
	return true
}
func (c *commandDefinition) validatePath(parameterName string, path string) bool {

	if parameterName == "" {
		colors.Default.Error.Println(fmt.Sprintf("%s directory is required", parameterName))
		color.Unset()
		return false
	}
	fileInfo, fileError := os.Stat(c.source)
	if !c.validateFileInfo(parameterName, fileInfo, &fileError) {
		colors.Default.Error.Println(fmt.Sprintf("Error while validating parameter: %s", parameterName))
		color.Unset()
		return false
	}
	colors.Default.Success.Println(fmt.Sprintf("Parameter %s is valid", parameterName))
	color.Unset()
	return true
}
func (c *commandDefinition) validateFlags() bool {

	spinner := helpers.NewSpinner()
	spinner.SetMessage("Validating flags...")
	spinner.Start()
	if !c.validatePath("source", c.source) {
		return false
	}
	if !c.validatePath("destination", c.destination) {
		return false
	}
	spinner.Finish()
	colors.Default.Success.Println("Flags are valid")
	color.Unset()
	return true
}

func (c *commandDefinition) runCommand(cmd *cobra.Command, args []string) {
	if !c.validateFlags() {
		os.Exit(1)
	}
	spinner := helpers.NewSpinner()
	spinner.SetMessage(fmt.Sprintf("Copying files from %s to %s...", c.source, c.destination))
	spinner.Start()

	err := filepath.Walk(c.source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Calculate the relative path from the source directory
		relativePath, err := filepath.Rel(c.source, path)
		if err != nil {
			return err
		}

		// Remove the .template suffix from the filename
		newFileName := strings.Replace(filepath.Base(path), ".template", "", 1)

		// Build the destination file path
		destinationPath := filepath.Join(c.destination, filepath.Dir(relativePath), newFileName)

		// Ensure the destination directory exists
		destDir := filepath.Dir(destinationPath)
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			spinner.Finish()
			colors.Default.Error.Printf("Error creating directory %s: %v\n", destDir, err)
			color.Unset()
			os.Exit(1)
		}

		// Copy the file
		spinner.SetMessage(fmt.Sprintf("Copying file %s to %s", path, destinationPath))
		err = helpers.CopyFile(path, destinationPath)
		if err != nil {
			spinner.Finish()
			colors.Default.Error.Printf("Error copying file %s to %s: %v\n", path, destinationPath, err)
			color.Unset()
			os.Exit(1)
		}

		spinner.SetMessage(fmt.Sprintf("File %s copied to %s", path, destinationPath))
		return nil
	})

	spinner.Finish()
	if err != nil {
		colors.Default.Error.Println("Error during file copy operation")
		colors.Default.Error.Println(err)
		color.Unset()
		os.Exit(1)
	}

	colors.Default.Success.Println(fmt.Sprintf("Files copied successfully from %s to %s", c.source, c.destination))
	color.Unset()
}

func (c *commandDefinition) createCommandDefinition() *cobra.Command {
	command := &cobra.Command{
		Use:   "templates",
		Short: "Bootstrap templates utility",
		Long:  `Copies the variables templates from a directory, to a desired destination, removing .template from the file name`,
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

}
func (c *commandDefinition) defineFlags() {
	c.cmd.Flags().StringVar(&c.source, "src", "", "Source directory")
	c.cmd.Flags().StringVar(&c.destination, "dest", "", "Destination directory")
	c.cmd.MarkFlagRequired("src")
	c.cmd.MarkFlagRequired("dest")
	c.cmd.MarkFlagsRequiredTogether("src", "dest")

}
func CreateCommand(context *context.NiOsContext, parentCommand *cobra.Command) *cobra.Command {
	c := commandDefinition{
		context:       context,
		parentCommand: parentCommand,
		source:        "",
		destination:   "",
	}
	return c.createCommandDefinition()
}
