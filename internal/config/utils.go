package config

import (
	"os"

	"github.com/spf13/cobra"
)

func GetUserHomeDir() string {
	// Search for a config file in default locations.
	home, err := os.UserHomeDir()
	// Only panic if we can't get the home directory.
	cobra.CheckErr(err)
	return home
}

func GetConfigPath(userHomeDir string) string {
	return userHomeDir + "/" + cfgPath
}
