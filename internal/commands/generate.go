package commands

import "github.com/spf13/cobra"

func cmdGenerate() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate [command]",
		Short: "Generates resources. See subcommands",
	}
	kubeConfigCmdOptions := newKubeconfigOptions()
	generateCmd.AddCommand(cmdKubeconfig(kubeConfigCmdOptions))
	return generateCmd
}
