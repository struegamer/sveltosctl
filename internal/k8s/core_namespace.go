package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/projectsveltos/sveltosctl/internal/logging"
)

type CoreClient struct {
	client    client.Client
	clientSet kubernetes.Interface
	logger    logging.Logger
	verbose   bool
}

func NewCoreClient(client *Client) *CoreClient {
	return NewCoreClientWithLogger(client, logging.NewKlogLogger(nil))
}

func NewCoreClientWithLogger(client *Client, logger logging.Logger) *CoreClient {
	return &CoreClient{
		client:    client.GetK8SClient(),
		clientSet: client.GetK8SClientSet(),
		logger:    logger,
		verbose:   false,
	}
}

// List fetches all namespaces (no context)
func (c *CoreClient) List() (*corev1.NamespaceList, error) {
	return c.ListCtx(context.TODO())
}

// ListCtx fetches all namespaces (with context)
func (c *CoreClient) ListCtx(ctx context.Context) (*corev1.NamespaceList, error) {

	c.logger.Info("Get all CoreResources")

	list := &corev1.NamespaceList{}
	err := c.client.List(ctx, list, &client.ListOptions{})
	return list, err
}

// Create namespace (without context)
func (c *CoreClient) Create(namespace string) error {
	return c.CreateCtx(context.TODO(), namespace)
}

// CreateCtx namespace (with context)
func (c *CoreClient) CreateCtx(ctx context.Context, name string) error {
	currentNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	err := c.client.Create(ctx, currentNs)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		c.logger.Info(fmt.Sprintf("Failed to create Namespace %s: %v",
			name, err))
		return err
	}
	return nil
}

// Get fetches one specified namespace (without context)
func (c *CoreClient) Get(name string) (*corev1.Namespace, error) {
	return c.GetCtx(context.TODO(), name)
}

// GetCtx fetches one specified namespace (with context)
func (c *CoreClient) GetCtx(ctx context.Context, name string) (*corev1.Namespace, error) {
	currentNs := &corev1.Namespace{}
	if err := c.client.Get(ctx, types.NamespacedName{Name: name}, currentNs); err != nil {
		return nil, err
	}
	return currentNs, nil
}

// Exists checks if namespace exists in cluster
func (c *CoreClient) Exists(name string) bool {
	_, err := c.Get(name)
	return err == nil
}
