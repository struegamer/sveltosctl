package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	ictx "github.com/projectsveltos/sveltosctl/internal/ctx"
)

type clusterOptions struct {
	Namespace           string
	Cluster             string
	kubeconfig          string
	fleetClusterContext string
	pullMode            bool
	labels              []string
	serviceAccountToken bool
}

const (
	clusterOptionsContextKey = ictx.ContextKey("cmdClusterOptions")
)

func newClusterOptions() *clusterOptions {
	return &clusterOptions{}
}

func addClusterOptionsToContext(ctx context.Context, cmdClusterOptions *clusterOptions) context.Context {
	ctx = context.WithValue(ctx, clusterOptionsContextKey, cmdClusterOptions)
	return ctx
}

func getClusterOptions(ctx context.Context) *clusterOptions {
	kco, ok := ctx.Value(clusterOptionsContextKey).(*clusterOptions)
	if !ok || kco == nil {
		panic("cluster options not found in context")
	}
	return kco
}

func clusterCmd(cmdClusterOptions *clusterOptions) *cobra.Command {
	checkFlags := func(cmd *cobra.Command) error {
		kubeCfg, err := cmd.Flags().GetString("kubeconfig")
		if err != nil {
			return err
		}
		fleetClusterContextName, err := cmd.Flags().GetString("fleet-cluster-context")
		if err != nil {
			return err
		}
		if kubeCfg == "" && fleetClusterContextName == "" {
			return errors.New("no --kubeconfig or --fleet-cluster-context specified, provide at least one of them")
		}
		cluster, err := cmd.Flags().GetString("cluster")
		if err != nil {
			return err
		}
		if cluster == "" {
			return errors.New("cluster name must be specified")
		}
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}
		if namespace == "" {
			return errors.New("namespace name must be specified")
		}
		return nil
	}
	cmdCluster := &cobra.Command{
		Use:   "cluster",
		Short: "The register cluster command registers a cluster to be managed by Sveltos.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ctx = addClusterOptionsToContext(ctx, cmdClusterOptions)
			cmd.SetContext(ctx)
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			clusterCmdOptions := getClusterOptions(cmd.Context())
			err := checkFlags(cmd)
			if err != nil {
				return err
			}
			labels, err := stringToMap(clusterCmdOptions.labels)
			if err != nil {
				return err
			}
			cmd.Println(labels)
			cmd.Println(clusterCmdOptions)
			return nil
		},
	}

	cmdCluster.Flags().StringVar(&cmdClusterOptions.Namespace, "namespace", "", `Specifies the namespace where Sveltos will create a resource (SveltosCluster) to represent
                                         the registered cluster.`)
	cmdCluster.Flags().StringVar(&cmdClusterOptions.Cluster, "cluster", "", `Defines a name for the registered cluster within Sveltos.`)
	cmdCluster.Flags().StringVar(&cmdClusterOptions.kubeconfig, "kubeconfig", "",
		`Provides the path to a file containing the kubeconfig for the Kubernetes cluster
you want to register.
If you don't have a kubeconfig file yet, you can use the "sveltosctl generate kubeconfig" command.
Be sure to point that command to the specific cluster you want to manage.
This will help you create the necessary kubeconfig file before registering the cluster
\with Sveltos.
Either --kubeconfig or --fleet-cluster-context must be provided.`)
	cmdCluster.Flags().StringVar(&cmdClusterOptions.fleetClusterContext, "fleet-cluster-context", "",
		`If your kubeconfig has multiple contexts:
- One context points to the management cluster (default one)
- Another context points to the cluster you actually want to manage;
In this case, you can specify the context name with the --fleet-cluster-context flag.
This tells the command to use the specific context to generate a Kubeconfig Sveltos
can use and then create a SveltosCluster with it so you don't have to provide kubeconfig
Either --kubeconfig or --fleet-cluster-context must be provided.`)
	cmdCluster.Flags().BoolVar(&cmdClusterOptions.pullMode, "pullmode", false,
		`this registers a cluster in pull mode. When enabled, the managed cluster will actively
fetch its configurations from the management cluster, which is ideal for scenarios with
firewall restrictions or when direct inbound access to the managed cluster is undesirable.
This flag outputs the specialized YAML configuration that needs to be applied to the managed
cluster to complete its setup.`)
	cmdCluster.Flags().StringSliceVar(&cmdClusterOptions.labels, "labels", []string{},
		`This option allows you to specify labels for the SveltosCluster resource
being created. The format for labels is <key1=value1,key2=value2>, where each key-value
pair is separated by a comma (,) and the key and value are separated by an equal sign (=).
You can define multiple labels by adding more key-value pairs separated by commas.`)
	return cmdCluster
}

func stringToMap(data []string) (map[string]string, error) {
	const keyValueLength = 2
	result := make(map[string]string)
	for kvPair := range data {
		kv := strings.Split(data[kvPair], "=")
		if len(kv) != keyValueLength {
			return nil, fmt.Errorf("invalid key-value pair format %q", kvPair)
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		result[key] = value
	}
	return result, nil
}
