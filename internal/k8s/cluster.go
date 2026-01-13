package k8s

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Cluster struct {
	client     *Client
	restConfig *rest.Config
	clientSet  *kubernetes.Clientset
	scheme     *runtime.Scheme
}

type ApiSchemaFunc func(s *runtime.Scheme) error

// NewCluster does all the heavy lifting for creating k8s client,configs,schemas etc.
func NewCluster(apiSchemas ...ApiSchemaFunc) *Cluster {
	cluster := &Cluster{
		scheme: runtime.NewScheme(),
	}
	if len(apiSchemas) != 0 {
		for _, apiSchema := range apiSchemas {
			if err := cluster.addSchema(apiSchema); err != nil {
				panic(err)
			}
		}
	}
	if err := cluster.initCluster(); err != nil {
		panic(err)
	}
	return cluster
}

func (cluster *Cluster) initCluster() error {
	restConfig, err := ctrl.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config %w", err)
	}
	restConfig.QPS = 100
	restConfig.Burst = 100
	cluster.restConfig = restConfig

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("error in getting access to K8S: %w", err)
	}
	cluster.clientSet = clientSet

	c, err := NewClient(restConfig, cluster.scheme)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	cluster.client = c
	return nil
}

func (cluster *Cluster) addSchema(apiSchema ApiSchemaFunc) error {
	if err := apiSchema(cluster.scheme); err != nil {
		return err
	}
	return nil
}
