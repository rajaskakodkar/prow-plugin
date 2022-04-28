package main

import (
	"context"
	"log"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"
)

var (
	namespace     = "prow"
	secretsFolder = ".secrets"

	clientset *kubernetes.Clientset

	// todo(knabben) - read from a config file
	hmacSecret           = filepath.Join(homedir.HomeDir(), secretsFolder, "hmac")
	githubSecret         = filepath.Join(homedir.HomeDir(), secretsFolder, "github")
	githuboaSecret       = filepath.Join(homedir.HomeDir(), secretsFolder, "githubOAuth")
	cookieSecret         = filepath.Join(homedir.HomeDir(), secretsFolder, "cookieSecret")
	serviceAccountSecret = filepath.Join(homedir.HomeDir(), secretsFolder, "serviceAccount")
)

func renderSecretSpec(name, namespace string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{},
	}
}

func createRequiredSecrets(kubeConfig string) error {
	log.Println("Creating secrets.")
	clientset = getClientSet(kubeConfig)

	secrets := []struct {
		name           string
		fileSources    []string
		literalSources []string
	}{
		/* kubectl create secret -n prow generic hmac-token --from-file= hmac= /home/ubuntu/prow/hook-secret */
		{
			name:           "hmac-token",
			fileSources:    []string{"hmac=" + hmacSecret},
			literalSources: []string{},
		},
		/* kubectl create secret -n prow generic github-token --from-file=cert=/home/ubuntu/prow/vagator.pem --from-literal=appid="153819" */
		{
			name:           "github-token",
			fileSources:    []string{"cert=" + githubSecret},
			literalSources: []string{"appid=\"153819\""},
		},
		/*kubectl create secret generic github-oauth-config --from-file=secret=github-oauth-config.yaml -n prow*/
		{
			name:           "github-oauth-config",
			fileSources:    []string{"secret=" + githuboaSecret},
			literalSources: []string{},
		},
		{
			name:           "cookie",
			fileSources:    []string{"secret=" + cookieSecret},
			literalSources: []string{},
		},
		{
			name:           "gcs-credentials",
			fileSources:    []string{serviceAccountSecret},
			literalSources: []string{},
		},
	}
	for _, secret := range secrets {
		log.Printf("Creating a secret: %s", secret.name)
		if err := CreateNewSecret(secret.name, secret.fileSources, secret.literalSources); err != nil && !strings.Contains(err.Error(), "already exists") {
			return err
		}
	}
	return nil
}

func CreateNewSecret(name string, fileSources, literalSources []string) error {
	secretsClient := clientset.CoreV1().Secrets(namespace)

	secret := renderSecretSpec(name, namespace)
	if err := createSecretsFileSource(secret, fileSources); err != nil {
		return err
	}
	if err := createSecretsliteralSource(secret, literalSources); err != nil {
		return err
	}
	_, err := secretsClient.Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// all good
	//fmt.Println(result)
	return nil
}

// createSecretsLiteralSources add new keyNames from Literal Sources directly
func createSecretsliteralSource(secret *corev1.Secret, literalSources []string) error {
	for _, literalSource := range literalSources {
		keyName, value, err := ParseLiteralSource(literalSource)
		if err != nil {
			return err
		}
		if err = addKeyFromLiteralToSecret(secret, keyName, []byte(value)); err != nil {
			return err
		}
	}
	return nil
}

// createSecretsFileSource add new keyNames from file Sources directly
func createSecretsFileSource(secret *corev1.Secret, fileSources []string) error {
	for _, fileSource := range fileSources {
		keyName, filePath, err := ParseFileSource(fileSource)
		if err != nil {
			return err
		}

		if err := addKeyFromFileToSecret(secret, keyName, filePath); err != nil {
			return err
		}
	}
	return nil
}
