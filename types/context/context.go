package context

import "github.com/spf13/cobra"

type CobraContext struct {
	ConfigurationFilePath string
}
type NiOsContext struct {
	NiOsCmd                      cobra.Command
	NixConfigurationsDir         string
	CreateVariablesFromTemplates bool
	SetupDisko                   bool
	DiskoConfigPath              string
	CreateHardwareConfiguration  bool
}
