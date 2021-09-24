package controllers

import (
	"context"
	"errors"
	"fmt"
	_ "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var (
	clientSet *kubernetes.Clientset
	err error
)

//New initializes the client with the given namespace and secret
func New(isLocal bool) (error) {
	//to-do move clientset initization to common layer
	var config *rest.Config
	if isLocal {
		home, exists := os.LookupEnv("HOME")
		if !exists {
			home = "/root"
		}
		configPath := filepath.Join(home, ".kube", "config")
		config, _ = clientcmd.BuildConfigFromFlags("", configPath)
	} else  {
		config, _ = rest.InClusterConfig()
	}

	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		return errors.New("Failed to create kubernetes client.  " +err.Error())
	}
	return nil
}

//UpdateSecret updates the secret with the key value for the given secret in the namespace
func UpdateSecret( namespace , secretName , key, value string ) ( error) {
	secretsClient := clientSet.CoreV1().Secrets(namespace)
	sec,_ := GetSecret(namespace,secretName)
	fmt.Println("Adding new key/value pair to secret as a string (StringData)")
	sec.Data[key] = []byte(value)

	_, err := secretsClient.Update(context.Background(),sec,metaV1.UpdateOptions{})

	if err != nil {
		return  err
	}

	return nil
}

//CreateSecret creates a new secret
func CreateSecret( namespace string , secret *v1.Secret) (error) {
	secretsClient := clientSet.CoreV1().Secrets(namespace)
	_, err := secretsClient.Create(context.Background(),secret,metaV1.CreateOptions{})

	if err != nil {
		return  err
	}
	return nil
}

//DeleteSecret deletes a secret with the given name
func DeleteSecret( namespace string , secretName string ) (error) {
	secretsClient := clientSet.CoreV1().Secrets(namespace)
	err := secretsClient.Delete(context.Background(),secretName,metaV1.DeleteOptions{})

	if err != nil {
		return err
	}
	return nil
}

// GetSecret gets a secret with the name
func GetSecret( namespace string , name string ) (*v1.Secret,error) {
	secretsClient := clientSet.CoreV1().Secrets(namespace)
	secret,err := secretsClient.Get(context.Background(),name,metaV1.GetOptions{})

	if err != nil {
		return nil,err
	}

	return secret,nil
}

