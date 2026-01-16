package k8s

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// CreateSecretCtx create a secret (with context)
func (c *CoreClient) CreateSecretCtx(ctx context.Context, namespace, saName string) error {
	c.logger.Info(fmt.Sprintf("Create Secret %s/%s", namespace, saName))
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

	err := c.client.Create(ctx, currentSecret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		c.logger.Info(fmt.Sprintf("Failed to create Secret %s/%s: %v",
			namespace, saName, err))
		return err
	}
	return nil
}

// CreateSecret create a secret (with context)
func (c *CoreClient) CreateSecret(namespace, saName string) error {
	return c.CreateSecretCtx(context.TODO(), namespace, saName)
}

func (c *CoreClient) GetSecretCtx(ctx context.Context, namespace, saName string) (*corev1.Secret, error) {
	c.logger.Info(fmt.Sprintf("Get Secret %s/%s", namespace, saName))
	currentSecret := &corev1.Secret{}
	err := c.client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: saName}, currentSecret)
	if err != nil {
		return nil, err
	}
	return currentSecret, nil
}
func (c *CoreClient) GetSecret(namespace, saName string) (*corev1.Secret, error) {
	return c.GetSecretCtx(context.Background(), namespace, saName)
}

// GetTokenCtx
func (c *CoreClient) GetTokenCtx(ctx context.Context, namespace, secretName string) (string, error) {
	getSecret := func() (string, error) {
		secret := &corev1.Secret{}
		err := c.client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: secretName}, secret)
		if err != nil {
			return "", err
		}
		if secret.Data == nil {
			return "", fmt.Errorf("secret %s/%s has no data", namespace, secretName)
		}
		v, ok := secret.Data["token"]
		if !ok {
			return "", fmt.Errorf("secret %s/%s has no token", namespace, secretName)
		}

		return string(v), nil
	}

	c.logger.Info(fmt.Sprintf("Get Token for  %s/%s", namespace, secretName))
	retries := 0
	const maxRetries = 5
	const maxBackoff = 32 * time.Second
	const baseTime = 1 * time.Second
	for {
		if retries >= maxRetries {
			return "", fmt.Errorf("retries exceeded")
		}
		backoff := time.Duration(math.Min(float64(baseTime)*math.Pow(2, float64(retries)), float64(maxBackoff)))
		jitter := time.Duration(rand.Float64() + float64(backoff)*0.5)
		nextBackoff := backoff + jitter

		token, err := getSecret()
		if err == nil {
			return token, nil
		}
		time.Sleep(nextBackoff)
		retries++
	}
}

// GetToken
func (c *CoreClient) GetToken(namespace, secretName string) (string, error) {
	return c.GetTokenCtx(context.TODO(), namespace, secretName)
}

// CreateTokenCtx
func (c *CoreClient) CreateTokenCtx(ctx context.Context, namespace, saName string, expirationSeconds int) (*authenticationv1.TokenRequest, error) {
	expiration := int64(expirationSeconds)
	treq := &authenticationv1.TokenRequest{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      saName,
		},
	}

	if expirationSeconds != 0 {
		treq.Spec = authenticationv1.TokenRequestSpec{
			ExpirationSeconds: &expiration,
		}
	}

	tokenRequest, err := c.clientSet.CoreV1().ServiceAccounts(namespace).CreateToken(ctx, saName, treq, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return tokenRequest, nil
}

// CreateToken
func (c *CoreClient) CreateToken(namespace, secretName string, expirationSeconds int) (*authenticationv1.TokenRequest, error) {
	return c.CreateTokenCtx(context.TODO(), namespace, secretName, expirationSeconds)
}
