package nioscli

import (
	"github.com/dev-vinicius-andrade/nioscli/types/context"
)

func createContext() *context.NiOsContext {
	//defaultDotfilesDir := filepath.Join(homeDir, constants.DefaultToolDirName)
	//defaultConfigurationsFolderPath := filepath.Join(defaultDotfilesDir, constants.DefaultToolDirName)
	//defaultEnvironmentVariablesFilePath := filepath.Join(defaultConfigurationsFolderPath, constants.EnvironmentVariablesFileName)
	return &context.NiOsContext{
		// HomeDir:             homeDir,
		// DotfilesDir:         "",
		// DotfilesInformation: information,
		// Cobra: context.CobraContext{
		// 	ConfigurationFilePath: "",
		// },
		// ConfigurationsFolder:         defaultConfigurationsFolderPath,
		// EnvironmentVariablesFilePath: defaultEnvironmentVariablesFilePath,
	}
}
