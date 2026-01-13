package k8s

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	logs "github.com/projectsveltos/libsveltos/lib/logsettings"
	"github.com/projectsveltos/sveltosctl/internal/logging"
)

type NamespaceClient struct {
	client client.Client
	logger logr.Logger
}

func NewNamespaceClient(client client.Client) *NamespaceClient {
	return NewNamespaceClientWithLogger(client, logging.InitLogger())
}

func NewNamespaceClientWithLogger(client client.Client, logger logr.Logger) *NamespaceClient {
	return &NamespaceClient{
		client: client,
		logger: logger,
	}
}

// List fetches all namespaces (no context)
func (c *NamespaceClient) List() (*corev1.NamespaceList, error) {
	return c.ListCtx(context.TODO())
}

// ListCtx fetches all namespaces (with context)
func (c *NamespaceClient) ListCtx(ctx context.Context) (*corev1.NamespaceList, error) {
	c.logger.V(logs.LogDebug).Info("Get all Namespaces")
	list := &corev1.NamespaceList{}
	err := c.client.List(ctx, list, &client.ListOptions{})
	return list, err
}

// Create namespace (without context)
func (c *NamespaceClient) Create(namespace string) error {
	return c.CreateCtx(context.TODO(), namespace)
}

// CreateCtx namespace (with context)
func (c *NamespaceClient) CreateCtx(ctx context.Context, name string) error {
	currentNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	err := c.client.Create(ctx, currentNs)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		c.logger.V(logs.LogDebug).Info(fmt.Sprintf("Failed to create Namespace %s: %v",
			name, err))
		return err
	}
	return nil
}

// Get fetches one specified namespace (without context)
func (c *NamespaceClient) Get(name string) (*corev1.Namespace, error) {
	return c.GetCtx(context.TODO(), name)
}

// GetCtx fetches one specified namespace (with context)
func (c *NamespaceClient) GetCtx(ctx context.Context, name string) (*corev1.Namespace, error) {
	currentNs := &corev1.Namespace{}
	if err := c.client.Get(ctx, types.NamespacedName{Name: name}, currentNs); err != nil {
		return nil, err
	}
	return currentNs, nil
}
