package hosts

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/dev-vinicius-andrade/nioscli/helpers"
	"github.com/dev-vinicius-andrade/nioscli/types/colors"
	"github.com/dev-vinicius-andrade/nioscli/types/context"
	"github.com/spf13/cobra"
)

type commandDefinition struct {
	cmd                               *cobra.Command
	context                           *context.NiOsContext
	parentCommand                     *cobra.Command
	path                              string
	name                              string
	templatesPath                     string
	hostName                          string
	useScopedCommonVars               bool
	commonVarsTemplateSearchText      string
	commonVarsTemplateReplaceText     string
	hostVariablesTemplateFileName     string
	commonVariablesTemplateFileName   string
	hostConfigurationTemplateFileName string
	hostConfigurationFileName         string
	hostVariablesFileName             string
	commonVariablesFileName           string
}

func (c *commandDefinition) handleUseScopedCommonVars(hostConfigurationFilePath string) {
	if !c.useScopedCommonVars {
		return
	}
	content, err := os.ReadFile(hostConfigurationFilePath)
	if err != nil {
		colors.Default.Error.Println(fmt.Errorf("error opening file: %v", err))
		os.Exit(1)
	}
	contentStr := string(content)
	contentStr = strings.Replace(contentStr, c.commonVarsTemplateSearchText, c.commonVarsTemplateReplaceText, 1)
	err = os.WriteFile(hostConfigurationFilePath, []byte(contentStr), 0644)
	if err != nil {
		colors.Default.Error.Println(fmt.Errorf("error writing file: %v", err))
		return
	}
	colors.Default.Info.Printf("Host configuration file %s is now using scoped common vars.\n", hostConfigurationFilePath)
}
func (c *commandDefinition) copyVariables(hostConfigurationPath, templatesPath string) {
	variablesPath := path.Join(templatesPath, "variables")
	hostVariablesPath := path.Join(hostConfigurationPath, "variables")
	os.MkdirAll(hostVariablesPath, os.ModePerm)

	hostVariablesTemplateFilePath := path.Join(variablesPath, c.hostVariablesTemplateFileName)
	commonVariablesTemplateFilePath := path.Join(variablesPath, c.commonVariablesTemplateFileName)
	hostVariablesFilePath := path.Join(hostVariablesPath, c.hostVariablesFileName)
	helpers.CopyFile(hostVariablesTemplateFilePath, hostVariablesFilePath)
	if c.useScopedCommonVars {
		helpers.CopyFile(commonVariablesTemplateFilePath, path.Join(hostVariablesPath, c.commonVariablesFileName))
		colors.Default.Info.Printf("Common Scoped variables copied to %s\n", hostVariablesPath)
		return
	}
	colors.Default.Info.Printf("Variables copied to %s\n", hostVariablesPath)
}
func (c *commandDefinition) runCommand(cmd *cobra.Command, args []string) {
	hostsConfigurationPath := path.Join(c.path, "hosts")
	hostConfigurationPath := path.Join(hostsConfigurationPath, c.name)

	hostConfigurationFilePath := path.Join(hostConfigurationPath, c.hostConfigurationFileName)

	templatesPath := path.Join(c.path, helpers.TernaryString(strings.HasPrefix(c.templatesPath, "/"), c.templatesPath[1:], c.templatesPath))
	hostConfigurationTemplateFilePath := path.Join(templatesPath, c.hostConfigurationTemplateFileName)
	os.MkdirAll(hostConfigurationPath, os.ModePerm)
	helpers.CopyFile(hostConfigurationTemplateFilePath, hostConfigurationFilePath)
	c.handleUseScopedCommonVars(hostConfigurationFilePath)
	c.copyVariables(hostConfigurationPath, templatesPath)
}
func (c *commandDefinition) createCommandDefinition() *cobra.Command {
	command := &cobra.Command{
		Use:   "hosts",
		Short: "A host configuration template utility",
		Long:  "A host configuration template utility",
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
	// Add subcommands here
}
func (c *commandDefinition) defineFlags() {
	// Add flags here

	c.cmd.Flags().StringVar(&c.path, "path", "", `
		Path to the nixos configurations directory
		------------------------------
		This path will be used to copy the templates to the new host configuration
	`)
	c.cmd.Flags().StringVar(&c.name, "name", "", `
		Name of the host configuration
		------------------------------
		This name will be used to create a directory inside the hosts folder
	`)
	c.cmd.Flags().StringVar(&c.templatesPath, "templates-path", "templates", `
		 Path to the templates directory
		------------------------------
		This path will be used to copy the templates to the new host configuration
		If not provided it will use the path flag and append /templates

	`)
	c.cmd.Flags().StringVar(&c.hostName, "host-name", "", `
		Hostname of the host
		------------------------------
		This hostname will be use to overwrite the variable host.hostname in the configuration.nix file.
		If not provided, it will generate a random hostname
	`)
	c.cmd.Flags().BoolVar(&c.useScopedCommonVars, "use-scoped-common-vars", false, `
		Use scoped common vars
		------------------------------
		When this flag is set, the common vars will be copied to the host configuration directory
		Then it will modify the configuration.nix of the host to use the common vars from the host configuration directory	
	`)
	c.cmd.Flags().StringVar(&c.commonVarsTemplateSearchText, "common-vars-template-search-text", "common_vars = import ../../variables/common.nix;", `
		Common vars template search text
		------------------------------
		This text will be used to search the common vars import in the configuration.template.nix file
	`)
	c.cmd.Flags().StringVar(&c.commonVarsTemplateReplaceText, "common-vars-template-replace-text", "common_vars = import ./variables/common.nix;", `
		Common vars template replace text
		------------------------------
		This text will be used to replace the common vars import in the configuration.template.nix file
	`)
	c.cmd.Flags().StringVar(&c.hostVariablesTemplateFileName, "host-variables-template-file-name", "host.template.nix", `
		Host variables template file name
		------------------------------
		The name of the host variables template file
	`)
	c.cmd.Flags().StringVar(&c.commonVariablesTemplateFileName, "common-variables-template-file-name", "common.template.nix", `
		Common variables template file name
		------------------------------
		The name of the common variables template file
	`)
	c.cmd.Flags().StringVar(&c.hostConfigurationTemplateFileName, "host-configuration-template-file-name", "configuration.template.nix", `
		Host configuration template file name
		------------------------------
		The name of the host configuration template file
	`)
	c.cmd.Flags().StringVar(&c.hostConfigurationFileName, "host-configuration-file-name", "configuration.nix", `
		Host configuration file name
		------------------------------
		The name of the host configuration file
	`)

	c.cmd.Flags().StringVar(&c.hostVariablesFileName, "host-variables-file-name", "host.nix", `
		Host variables file name
		------------------------------
		The name of the host variables file
	`)
	c.cmd.Flags().StringVar(&c.commonVariablesFileName, "common-variables-file-name", "common.nix", `
		Common variables file name
		------------------------------
		The name of the common variables file
	`)
	c.cmd.MarkFlagRequired("path")
	c.cmd.MarkFlagRequired("name")
}
func CreateCommand(context *context.NiOsContext, parentCommand *cobra.Command) *cobra.Command {
	command := &commandDefinition{
		context:       context,
		parentCommand: parentCommand,
		//Add your command definition extra default values here
	}
	return command.createCommandDefinition()
}
