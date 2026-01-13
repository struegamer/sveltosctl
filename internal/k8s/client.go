package k8s

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	client client.Client
}

func NewClient(config *rest.Config, scheme *runtime.Scheme) (*Client, error) {
	var err error
	k8sClient := &Client{}
	k8sClient.client, err = createK8sClient(config, scheme)
	if err != nil {
		return nil, err
	}
	return k8sClient, nil
}

func (c *Client) SetK8sClient(client client.Client) {
	c.client = client
}

func createK8sClient(restConfig *rest.Config, scheme *runtime.Scheme) (client.Client, error) {
	c, err := client.New(restConfig, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	return c, nil
}
