package kubeconfig

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	logs "github.com/projectsveltos/libsveltos/lib/logsettings"
)

const (
	Projectsveltos = "projectsveltos"
)

func GenerateKubeconfigForServiceAccount(ctx context.Context, remoteRestConfig *rest.Config, namespace, serviceAccountName string, expirationSeconds int, create, display, satoken bool, logger logr.Logger) error {

	s := runtime.NewScheme()
	err := clientgoscheme.AddToScheme(s)
	if err != nil {
		return err
	}

	var remoteClient client.Client
	remoteClient, err = client.New(remoteRestConfig, client.Options{Scheme: s})
	if err != nil {
		return err
	}

	if create {
		err = createNamespace(ctx, remoteClient, namespace, logger)
		if err != nil {
			return err
		}
		err = createServiceAccount(ctx, remoteClient, namespace, serviceAccountName, logger)
		if err != nil {
			return err
		}
		err = createClusterRole(ctx, remoteClient, Projectsveltos, logger)
		if err != nil {
			return err
		}
		err = createClusterRoleBinding(ctx, remoteClient, Projectsveltos, Projectsveltos, namespace,
			serviceAccountName, logger)
		if err != nil {
			return err
		}
	} else {
		err = getNamespace(ctx, remoteClient, namespace)
		if err != nil {
			return err
		}
		err = getServiceAccount(ctx, remoteClient, namespace, serviceAccountName)
		if err != nil {
			return err
		}
	}

	var token string
	if satoken {
		if err := createSecret(ctx, remoteClient, namespace, serviceAccountName, logger); err != nil {
			return err
		}
		var err error
		token, err = getToken(ctx, remoteClient, namespace, serviceAccountName)
		if err != nil {
			return err
		}
	} else {
		tokenRequest, err := getServiceAccountTokenRequest(ctx, remoteRestConfig, namespace, serviceAccountName,
			expirationSeconds, logger)
		if err != nil {
			return err
		}
		token = tokenRequest.Token
	}

	logger.V(logs.LogDebug).Info("Get Kubeconfig from TokenRequest")
	data := getKubeconfigFromToken(remoteRestConfig, namespace, serviceAccountName, token)
	if display {
		//nolint: forbidigo // print kubeconfig
		fmt.Println(data)
	}

	return nil
}

func createSecret(ctx context.Context, c client.Client, namespace, saName string,
	logger logr.Logger) error {

	logger.V(logs.LogInfo).Info(fmt.Sprintf("Create Secret %s/%s", namespace, saName))
	currentSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      saName,
			Annotations: map[string]string{
				corev1.ServiceAccountNameKey: saName,
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}

	err := c.Create(ctx, currentSecret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		logger.V(logs.LogInfo).Info(fmt.Sprintf("Failed to create Secret %s/%s: %v",
			namespace, saName, err))
		return err
	}

	return nil
}

func getToken(ctx context.Context, c client.Client, namespace, secretName string) (string, error) {
	retries := 0
	const maxRetries = 5
	for {
		secret := &corev1.Secret{}
		err := c.Get(ctx, types.NamespacedName{Namespace: namespace, Name: secretName},
			secret)
		if err != nil {
			if retries < maxRetries {
				time.Sleep(time.Second)
				continue
			}
			return "", err
		}

		if secret.Data == nil {
			time.Sleep(time.Second)
			continue
		}

		v, ok := secret.Data["token"]
		if !ok {
			time.Sleep(time.Second)
			continue
		}

		return string(v), nil
	}
}

func getNamespace(ctx context.Context, remoteClient client.Client, name string) error {
	currentNs := &corev1.Namespace{}
	return remoteClient.Get(ctx, types.NamespacedName{Name: name}, currentNs)
}

func createNamespace(ctx context.Context, remoteClient client.Client, name string, logger logr.Logger) error {
	logger.V(logs.LogDebug).Info(fmt.Sprintf("Create namespace %s", name))
	currentNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	err := remoteClient.Create(ctx, currentNs)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		logger.V(logs.LogDebug).Info(fmt.Sprintf("Failed to create Namespace %s: %v",
			name, err))
		return err
	}

	return nil
}

func getServiceAccount(ctx context.Context, remoteClient client.Client, namespace, name string) error {
	currentSA := &corev1.ServiceAccount{}
	return remoteClient.Get(ctx,
		types.NamespacedName{Namespace: namespace, Name: name},
		currentSA)
}

func createServiceAccount(ctx context.Context, remoteClient client.Client, namespace, name string,
	logger logr.Logger) error {

	logger.V(logs.LogDebug).Info(fmt.Sprintf("Create serviceAccount %s/%s", namespace, name))
	currentSA := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}

	err := remoteClient.Create(ctx, currentSA)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		logger.V(logs.LogDebug).Info(fmt.Sprintf("Failed to create ServiceAccount %s/%s: %v",
			namespace, name, err))
		return err
	}

	return nil
}

func createClusterRole(ctx context.Context, remoteClient client.Client, clusterRoleName string,
	logger logr.Logger) error {

	logger.V(logs.LogDebug).Info(fmt.Sprintf("Create ClusterRole %s", clusterRoleName))
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

	err := remoteClient.Create(ctx, clusterrole)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		logger.V(logs.LogDebug).Info(fmt.Sprintf("Failed to create ClusterRole %s: %v",
			clusterRoleName, err))
		return err
	}

	return nil
}

func createClusterRoleBinding(ctx context.Context, remoteClient client.Client,
	clusterRoleName, clusterRoleBindingName, serviceAccountNamespace, serviceAccountName string,
	logger logr.Logger) error {

	logger.V(logs.LogDebug).Info(fmt.Sprintf("Create ClusterRoleBinding %s", clusterRoleBindingName))
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
	err := remoteClient.Create(ctx, clusterrolebinding)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		logger.V(logs.LogDebug).Info(fmt.Sprintf("Failed to create clusterrolebinding %s: %v",
			clusterRoleBindingName, err))
		return err
	}

	return nil
}

// getServiceAccountTokenRequest returns token for a serviceaccount
func getServiceAccountTokenRequest(ctx context.Context, remoteRestConfig *rest.Config, serviceAccountNamespace, serviceAccountName string,
	expirationSeconds int, logger logr.Logger) (*authenticationv1.TokenRequestStatus, error) {

	expiration := int64(expirationSeconds)

	treq := &authenticationv1.TokenRequest{}

	if expirationSeconds != 0 {
		treq.Spec = authenticationv1.TokenRequestSpec{
			ExpirationSeconds: &expiration,
		}
	}

	clientset, err := kubernetes.NewForConfig(remoteRestConfig)
	if err != nil {
		return nil, err
	}

	logger.V(logs.LogDebug).Info(
		fmt.Sprintf("Create Token for ServiceAccount %s/%s", serviceAccountNamespace, serviceAccountName))
	var tokenRequest *authenticationv1.TokenRequest
	tokenRequest, err = clientset.CoreV1().ServiceAccounts(serviceAccountNamespace).
		CreateToken(ctx, serviceAccountName, treq, metav1.CreateOptions{})
	if err != nil {
		logger.V(logs.LogDebug).Info(
			fmt.Sprintf("Failed to create token for ServiceAccount %s/%s: %v",
				serviceAccountNamespace, serviceAccountName, err))
		return nil, err
	}

	return &tokenRequest.Status, nil
}

// getKubeconfigFromToken returns Kubeconfig to access management cluster from token.
func getKubeconfigFromToken(remoteRestConfig *rest.Config, namespace, serviceAccountName, token string) string {
	template := `apiVersion: v1
kind: Config
clusters:
- name: local
  cluster:
    server: %s
    certificate-authority-data: "%s"
users:
- name: %s
  user:
    token: %s
contexts:
- name: sveltos-context
  context:
    cluster: local
    namespace: %s
    user: %s
current-context: sveltos-context`

	data := fmt.Sprintf(template, remoteRestConfig.Host,
		base64.StdEncoding.EncodeToString(remoteRestConfig.CAData), serviceAccountName, token, namespace, serviceAccountName)

	return data
}
