package config

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/dev-vinicius-andrade/nioscli/helpers"
	"github.com/dev-vinicius-andrade/nioscli/types"
	"github.com/dev-vinicius-andrade/nioscli/types/colors"
	"github.com/dev-vinicius-andrade/nioscli/types/context"
	"github.com/dev-vinicius-andrade/nioscli/types/enums"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type commandDefinition struct {
	cmd           *cobra.Command
	context       *context.NiOsContext
	parentCommand *cobra.Command
	destination   string
	interactive   bool
	config        *types.NixDiskDevices
}

// func (c *commandDefinition) generateConfig(cmd *cobra.Command, args []string) {
// 	// get current os
// 	spinner := helpers.NewSpinner()
// 	spinner.SetMessage("Generating hardware configuration...")
// 	spinner.Start()
// 	cmdArgsBuilder := fmt.Sprintf("--root %s", c.root)
// 	if c.noFileSytem {
// 		cmdArgsBuilder += " --no-filesystem"
// 	}
// 	spinner.SetMessage(fmt.Sprintf("Running command:  %s", cmdArgsBuilder))

// 	commandString := fmt.Sprintf("nixos-generate-config %s", cmdArgsBuilder)

// 	command := exec.Command("/bin/sh", "-c", commandString)
// 	_, err := command.Output()
// 	spinner.Finish()
// 	if err != nil {
// 		colors.Default.Error.Println("Error generating hardware configuration")
// 		colors.Default.Error.Println(err)
// 		color.Unset()
// 		os.Exit(1)
// 	}
// 	colors.Default.Success.Println("Hardware configuration generated")
// 	color.Unset()
// }
// func (c *commandDefinition) moveHardwareConfigToDestination(cmd *cobra.Command, args []string) {
// 	if !c.moveFile {
// 		return
// 	}
// 	// if c.destination == "" {
// 	// 	colors.Default.Error.Println("Destination path is required")
// 	// 	color.Unset()
// 	// 	os.Exit(1)
// 	// }
// 	info, destinationError := os.Stat(c.destination)
// 	if os.IsNotExist(destinationError) {
// 		fmt.Println("Destination path does not exist")
// 		return
// 	}
// 	if !info.IsDir() {
// 		fmt.Println("Destination path is not a directory")
// 		return
// 	}
// 	source := fmt.Sprintf("%s/etc/nixos/hardware-configuration.nix", c.root)
// 	commandString := fmt.Sprintf("mv %s %s", source, c.destination)
// 	command := exec.Command("/bin/sh", "-c", commandString)

// 	_, err := command.Output()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// }
func (c *commandDefinition) listDevices() []types.DiskDeviceInformation {
	commandString := "lsblk -d -o NAME,TYPE,SIZE | grep disk | awk '{print \"/dev/\" $1, $3}'"
	command := exec.Command("/bin/sh", "-c", commandString)
	output, err := command.Output()
	if err != nil {
		colors.Default.Error.Println("Error listing devices")
		colors.Default.Error.Println(err)
		color.Unset()
		return nil
	}
	deviceList := strings.Split(strings.TrimSpace(string(output)), "\n")
	var devices []types.DiskDeviceInformation
	for _, device := range deviceList {
		parts := strings.Fields(device)
		if len(parts) == 2 {
			devices = append(devices, types.DiskDeviceInformation{
				Device: parts[0],
				Size:   parts[1],
			})
		}
	}
	return devices
}
func (c *commandDefinition) runNonInteractive(cmd *cobra.Command, args []string) {
	colors.Default.Error.Println("For now, only interactive mode is available, please use the --interactive flag")
	color.Unset()
	os.Exit(1)
}
func (c *commandDefinition) runInteractive(cmd *cobra.Command, args []string) {
	device := c.promptDeviceName()
	device.Device = c.selectDevicePrompt()
	device.Type = c.selectDeviceTypePrompt()
}
func (c *commandDefinition) promptDeviceName() *types.NixDiskDevice {
	isValid := false
	deviceName := ""
	for !isValid {
		prompt := promptui.Prompt{
			Label:   "Enter the device name",
			Default: "main",
			Validate: func(input string) error {
				trimedInput := strings.TrimSpace(input)

				if trimedInput == "" {
					return errors.New("you need to enter a valid device name")
				}
				if strings.Contains(trimedInput, " ") {
					return errors.New("you can't use spaces in the device name")
				}
				allowedSpecialCharsExpression := regexp.MustCompile(`[^\w\s-_]`)
				if allowedSpecialCharsExpression.MatchString(trimedInput) {
					return errors.New("you can't use special characters rather than - and _")

				}
				return nil
			},
			Pointer: promptui.PipeCursor,
		}
		result, err := prompt.Run()
		if err != nil {
			colors.Default.Error.Println("error getting device name")
			colors.Default.Error.Println(err)
			color.Unset()
			os.Exit(1)
		}
		deviceName = result
		isValid = true
	}

	return &types.NixDiskDevice{
		Name: &deviceName,
	}
}
func (c *commandDefinition) listDevicesType() []enums.NixDiskDeviceType {
	return []enums.NixDiskDeviceType{
		enums.NixDiskDeviceTypes.Disk,
	}
}
func (c *commandDefinition) selectDeviceTypePrompt() *enums.NixDiskDeviceType {
	devicesTypes := c.listDevicesType()
	var selectedIndex = -1
	for selectedIndex < 0 {
		prompt := promptui.Select{
			Label: "Select which the type of your device",
			Items: devicesTypes,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   colors.Default.Blue.Sprintf("→ {{ . | white | underline }}"),
				Inactive: "{{.Device | white }}",
			},
		}
		index, _, err := prompt.Run()

		if err != nil {
			colors.Default.Error.Println("Error selecting device type")
			colors.Default.Error.Println(err)
			color.Unset()
			os.Exit(1)
		}
		selectedIndex = index
		if selectedIndex >= 0 {
			break
		} else {
			colors.Default.Error.Println("You need to select at least one type")
			color.Unset()
		}

	}
	return &devicesTypes[selectedIndex]
}
func (c *commandDefinition) selectDevicePrompt() *string {
	devices := c.listDevices()
	var selectedIndex = -1
	for selectedIndex < 0 {
		prompt := promptui.Select{
			Label: "Select which disk you want to use",
			Items: devices,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ .Device }}?",
				Active:   colors.Default.Blue.Sprintf("→ {{ .Device | white | underline }}"),
				Inactive: "{{.Device | white }}",
				Details:  colors.Default.Gray.Sprintf("\t\tSize: {{ .Size | underline }}"),
				FuncMap:  promptui.FuncMap,
			},
		}
		index, _, err := prompt.Run()

		if err != nil {
			colors.Default.Error.Println("Error selecting devices")
			colors.Default.Error.Println(err)
			color.Unset()
			os.Exit(1)
		}
		selectedIndex = index
		if selectedIndex >= 0 {
			break
		} else {
			colors.Default.Error.Println("You need to select at least one device")
			color.Unset()
		}

	}
	return &devices[selectedIndex].Device

}
func (c *commandDefinition) runCommand(cmd *cobra.Command, args []string) {
	if !c.interactive {
		c.runNonInteractive(cmd, args)
		return

	} else {
		c.runInteractive(cmd, args)
	}
	// c.generateConfig(cmd, args)
	// c.moveHardwareConfigToDestination(cmd, args)
}
func (c *commandDefinition) createCommandDefinition() *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "Tool to generate disko configuration file",
		Long:  `The goal of this tool is to help you to generate the disko configuration file, so you can use it to format and partition your disks, with a cli tool, also you can use the interactive mode to help you to generate the configuration file`,
		Run: func(cmd *cobra.Command, args []string) {
			c.runCommand(cmd, args)
		},
	}
	c.cmd = command
	c.defineFlags()
	c.defineSubCommands()
	helpers.CobraHelper.AddCommandToParent(command, c.parentCommand)
	return &c.context.NiOsCmd
}
func (c *commandDefinition) defineFlags() {
	c.cmd.Flags().BoolVar(&c.interactive, "interactive", false, "Flag to run the command in interactive mode")
}
func (c *commandDefinition) defineSubCommands() {
	// 	config.CreateCommand(c.context, c.cmd)
}
func CreateCommand(context *context.NiOsContext, parentCommand *cobra.Command) *cobra.Command {
	c := commandDefinition{
		context:       context,
		parentCommand: parentCommand,

		//path:          "./disko.nix",
	}
	return c.createCommandDefinition()
}
