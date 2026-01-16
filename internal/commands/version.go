package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/projectsveltos/sveltosctl/internal/config"
)

var (
	gitVersion string
	gitCommit  string
)

func cmdVersion() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of sveltosctl",
		Long:  "Print the version of sveltosctl",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ctlConfig := config.CtlConfigFromContext(ctx)
			if ctlConfig == nil {
				return errors.New("could not find sveltosctl config")
			}
			if !ctlConfig.NoConsoleOutput {
				cmd.Println("Client Version: ", gitVersion)
				cmd.Println("Git commit: ", gitCommit)
			} else {
				logger := ctlConfig.Logger()
				logger.Info(fmt.Sprintf("Git commit: %s", gitCommit))
				logger.Info(fmt.Sprintf("Client Version: %s", gitVersion))
			}
			return nil
		},
	}
	return versionCmd
}
