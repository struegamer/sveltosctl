package config

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const envPrefix = "SVELTOSCTL"
const cfgPath = ".sveltosctl"

func (cfg *CtlConfig) Initialize(cmd *cobra.Command) error {
	if cfg.Verbose {
		cfg.logger.Logger().Info("Initializing Viper config")
	}
	viper.SetEnvPrefix(envPrefix)
	// Allow for nested keys in environment variables (e.g. `MYAPP_DATABASE_HOST`)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	viper.AutomaticEnv()
	if cfg.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfg.CfgFile)
	} else {
		// Search for a config file in default locations and panic if it can not be found
		home := GetUserHomeDir()

		// Search for a config file with the name "config" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/" + cfgPath)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	return nil
}
