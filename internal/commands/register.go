package commands

import "github.com/spf13/cobra"

func cmdRegister() *cobra.Command {
	registerCmd := &cobra.Command{
		Use:   "register [command]",
		Short: "Register a cluster",
	}
	clusterCmdOptions := newClusterOptions()
	registerCmd.AddCommand(clusterCmd(clusterCmdOptions))
	return registerCmd
}
