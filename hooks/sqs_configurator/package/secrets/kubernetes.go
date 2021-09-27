package secrets

import (
	"context"
	"errors"
	"go.uber.org/zap"
	_ "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var secretsClient coreV1Types.SecretInterface

// KubeSecret holds the handler for the secret and auxiliary info
type KubeSecret struct {
	secret    *v1.Secret
	namespace string
}

//New initializes the client with the given namespace and secret
func New(namespace string, secretName string, isLocal *bool) (KubeSecret, error) {
	var config *rest.Config
	if *isLocal {
		home, exists := os.LookupEnv("HOME")
		if !exists {
			home = "/root"
		}
		configPath := filepath.Join(home, ".kube", "config")
		config, _ = clientcmd.BuildConfigFromFlags("", configPath)
	} else {
		config, _ = rest.InClusterConfig()
	}

	clientSet, err := kubernetes.NewForConfig(config)

	if err != nil {
		zap.L().Debug("Failed to create K8s clientset")
		panic(err.Error())
	}
	secretsClient = clientSet.CoreV1().Secrets(namespace)
	secret, err := secretsClient.Get(context.Background(), secretName, metaV1.GetOptions{})

	if err != nil {
		return KubeSecret{}, err
	}

	sec := KubeSecret{
		secret:    secret,
		namespace: namespace,
	}

	return sec, nil

}

// UpdateSecret updates the secret in the given namespace with the values provided
func UpdateSecret(namespace, secretName, key, value string, isLocal *bool) error {
	sec, setupErr := New(namespace, secretName, isLocal)
	if setupErr != nil {
		zap.L().Info("Error creating the secret in the namespace")
		return errors.New("Error creating the secret in the namespace")
	}
	zap.L().Info("Adding new key/value pair to secret as a string (StringData)")
	sec.secret.Data[key] = []byte(value)

	_, err := secretsClient.Update(context.Background(), sec.secret, metaV1.UpdateOptions{})

	if err != nil {
		return err
	}

	return nil
}
