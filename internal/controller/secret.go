package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	thermopilotv1 "github.com/seipan/thermo-pilot-controller/api/v1"
)

type SwitchBotCredentials struct {
	Token  string
	Secret string
}

func GetSwitchBotCredentials(ctx context.Context, c client.Client, spec thermopilotv1.ThermoPilotSpec, namespace string) (*SwitchBotCredentials, error) {
	secretRef := spec.SecretRef
	tokenKey := secretRef.TokenKey
	if tokenKey == "" {
		tokenKey = "token"
	}
	secretKey := secretRef.SecretKey
	if secretKey == "" {
		secretKey = "secret"
	}
	secret := &corev1.Secret{}
	err := c.Get(ctx, types.NamespacedName{
		Name:      secretRef.Name,
		Namespace: namespace,
	}, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s: %w", secretRef.Name, err)
	}
	tokenBytes, exists := secret.Data[tokenKey]
	if !exists {
		return nil, fmt.Errorf("token key '%s' not found in secret %s", tokenKey, secretRef.Name)
	}
	secretBytes, exists := secret.Data[secretKey]
	if !exists {
		return nil, fmt.Errorf("secret key '%s' not found in secret %s", secretKey, secretRef.Name)
	}
	if len(tokenBytes) == 0 {
		return nil, fmt.Errorf("token value is empty in secret %s", secretRef.Name)
	}
	if len(secretBytes) == 0 {
		return nil, fmt.Errorf("secret value is empty in secret %s", secretRef.Name)
	}
	return &SwitchBotCredentials{
		Token:  string(tokenBytes),
		Secret: string(secretBytes),
	}, nil
}
