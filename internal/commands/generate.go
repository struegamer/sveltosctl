package commands

import (
	"context"

	"github.com/spf13/cobra"
)

func cmdGenerate(ctx context.Context) (*cobra.Command, context.Context) {
	ctx = newKubeConfigOptionsWithContext(ctx)
	generateCmd := &cobra.Command{
		Use:   "generate [command]",
		Short: "Generates resources. See subcommands",
	}
	generateCmd.AddCommand(cmdKubeconfig(ctx))
	return generateCmd, ctx
}
