package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ListServiceAccountsCtx list all serviceaccounts in namespace (with context)
func (c *CoreClient) ListServiceAccountsCtx(ctx context.Context, namespace string) (*corev1.ServiceAccountList, error) {
	if namespace == "" {
		namespace = "default"
	}
	list := &corev1.ServiceAccountList{}
	err := c.client.List(ctx, list, &client.ListOptions{Namespace: namespace})
	return list, err
}

// ListServiceAccounts list all serviceaccounts in namespace (without context)
func (c *CoreClient) ListServiceAccounts(namespace string) (*corev1.ServiceAccountList, error) {
	return c.ListServiceAccountsCtx(context.TODO(), namespace)
}

// CreateServiceAccountCtx creates a service account in namespace (with context)
func (c *CoreClient) CreateServiceAccountCtx(ctx context.Context, name, namespace string) error {
	srvAcc := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	err := c.client.Create(ctx, srvAcc)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		c.logger.Info(fmt.Sprintf("Failed to create ServiceAccount %s/%s: %v",
			namespace, name, err))
		return err
	}
	return nil
}

// CreateServiceAccount creates a service account in namespace (without context)
func (c *CoreClient) CreateServiceAccount(name, namespace string) error {
	return c.CreateServiceAccountCtx(context.TODO(), name, namespace)
}

// GetServiceAccountCtx retrieves service account in namespace (with context)
func (c *CoreClient) GetServiceAccountCtx(ctx context.Context, name, namespace string) (*corev1.ServiceAccount, error) {
	srvAcc := &corev1.ServiceAccount{}
	err := c.client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, srvAcc)
	if err != nil {
		return nil, err
	}
	return srvAcc, nil
}

// GetServiceAccount retrieves service account in namespace (without context)
func (c *CoreClient) GetServiceAccount(name, namespace string) (*corev1.ServiceAccount, error) {
	return c.GetServiceAccountCtx(context.TODO(), name, namespace)
}

// CreateClusterRoleCtx creates ClusterRole (with context)
func (c *CoreClient) CreateClusterRoleCtx(ctx context.Context, clusterRoleName string) error {
	c.logger.Info(fmt.Sprintf("Create ClusterRole %s", clusterRoleName))
	// Extends permission in addon-controller-role-extra
	clusterrole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"*"},
				APIGroups: []string{"*"},
				Resources: []string{"*"},
			},
			{
				Verbs:           []string{"*"},
				NonResourceURLs: []string{"*"},
			},
		},
	}
	err := c.client.Create(ctx, clusterrole)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		c.logger.Info(fmt.Sprintf("Failed to create ClusterRole %s: %v",
			clusterRoleName, err))
		return err
	}

	return nil
}

// CreateClusterRole creates ClusterRole (without context)
func (c *CoreClient) CreateClusterRole(clusterRoleName string) error {
	return c.CreateClusterRoleCtx(context.TODO(), clusterRoleName)
}

// CreateClusterRoleBindingCtx creates ClusterRoleBinding for clusterRole with serviceAccount in serviceAccount namespace (with context)
func (c *CoreClient) CreateClusterRoleBindingCtx(ctx context.Context, clusterRoleName, clusterRoleBindingName, serviceAccountNamespace, serviceAccountName string) error {
	c.logger.Info(fmt.Sprintf("Create ClusterRoleBinding %s", clusterRoleBindingName))
	clusterrolebinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleBindingName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
		},
		Subjects: []rbacv1.Subject{
			{
				Namespace: serviceAccountNamespace,
				Name:      serviceAccountName,
				Kind:      "ServiceAccount",
				APIGroup:  corev1.SchemeGroupVersion.Group,
			},
		},
	}
	err := c.client.Create(ctx, clusterrolebinding)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		c.logger.Info(fmt.Sprintf("Failed to create clusterrolebinding %s: %v",
			clusterRoleBindingName, err))
		return err
	}
	return nil
}

// CreateClusterRoleBinding creates ClusterRoleBinding for clusterRole with serviceAccount in serviceAccount namespace (without context)
func (c *CoreClient) CreateClusterRoleBinding(clusterRoleName, clusterRoleBindingName, serviceAccountNamespace, serviceAccountName string) error {
	return c.CreateClusterRoleBindingCtx(context.TODO(), clusterRoleName, clusterRoleBindingName, serviceAccountNamespace, serviceAccountName)
}
