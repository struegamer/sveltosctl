package commands

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	ictx "github.com/projectsveltos/sveltosctl/internal/ctx"

	"github.com/projectsveltos/sveltosctl/internal/config"
)

func RootCmd(ctlConfig *config.CtlConfig) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sveltosctl",
		Short: "CLI for sveltosctl",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := ctlConfig.Initialize(cmd)
			if err != nil {
				return err
			}
			ctx := cmd.Context()
			ctx = context.WithValue(ctx, ictx.CtlConfigCtxKey, ctlConfig)
			cmd.SetContext(ctx)
			return nil
		},
	}
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.PersistentFlags().StringVarP(&ctlConfig.CfgFile, "config-filename", "c", config.GetConfigPath(config.GetUserHomeDir())+"/config.yaml", "config file (default is config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&ctlConfig.NoConsoleOutput, "console-output", "o", false, "enable pretty console output, defaults to false")
	rootCmd.PersistentFlags().BoolVarP(&ctlConfig.Verbose, "verbose", "v", false, "enable verbose output, defaults to false")
	rootCmd.AddCommand(cmdVersion(), cmdGenerate(), cmdRegister())
	return rootCmd
}
