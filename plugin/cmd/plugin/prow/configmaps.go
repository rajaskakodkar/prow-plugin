package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"unicode/utf8"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
)

var (
	cmFolder    = ".secrets"
	configCM    = filepath.Join(homedir.HomeDir(), cmFolder, "config")
	pluginsCM   = filepath.Join(homedir.HomeDir(), cmFolder, "plugins")
	jobconfigCM = filepath.Join(homedir.HomeDir(), cmFolder, "job-config")
)

func renderConfigMapSpec(name, namespace string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func createRequiredConfigmap(kubeConfig string) error {
	log.Println("Creating ConfigMaps.")
	clientset = getClientSet(kubeConfig)

	configmaps := []struct {
		name        string
		fileSources []string
	}{
		/* kubectl create configmap config --from-file=config.yaml="${REPO_PATH}"/config/prow/config.yaml --dry-run=client -oyaml | kubectl apply -f - -n prow */
		{
			name:        "config",
			fileSources: []string{"config.yaml=" + configCM},
		},
		/* kubectl create configmap plugins --from-file=plugins.yaml="${REPO_PATH}"/config/prow/plugins.yaml --dry-run=client -oyaml | kubectl apply -f - -n prow */
		{
			name:        "plugins",
			fileSources: []string{"plugins.yaml=" + pluginsCM},
		},
		/* kubectl create configmap job-config --from-file=${JOB_CONFIG_PATH} --dry-run=client -oyaml | kubectl apply -f - -n prow */
		{
			name:        "job-config",
			fileSources: []string{jobconfigCM},
		},
	}
	for _, cm := range configmaps {
		log.Printf("Creating a config map: %s", cm.name)
		if err := CreateNewConfigMap(cm.name, cm.fileSources); err != nil && !strings.Contains(err.Error(), "already exists") {
			return err
		}
	}
	return nil

	fmt.Println(kubeConfig)
	return nil
}

func addKeyFromFileToConfigMap(configMap *corev1.ConfigMap, keyName, filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if utf8.Valid(data) {
		return addKeyFromLiteralToConfigMap(configMap, keyName, string(data))
	}
	// err = validateNewConfigMap(configMap, keyName)
	// if err != nil {
	// 	return err
	// }
	configMap.BinaryData[keyName] = data

	return nil
}

func addKeyFromLiteralToConfigMap(configMap *corev1.ConfigMap, keyName, data string) error {
	// err := validateNewConfigMap(configMap, keyName)
	// if err != nil {
	// 	return err
	// }
	configMap.Data[keyName] = data

	return nil
}

func CreateNewConfigMap(name string, fileSources []string) error {
	cmClient := clientset.CoreV1().ConfigMaps(namespace)

	configmap := renderConfigMapSpec(name, namespace)
	if err := createConfigMapFileSource(configmap, fileSources); err != nil {
		return err
	}
	_, err := cmClient.Create(context.TODO(), configmap, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// all good
	//fmt.Println(result)
	return nil
}

func createConfigMapFileSource(configmap *corev1.ConfigMap, fileSources []string) error {
	for _, fileSource := range fileSources {
		keyName, filePath, err := ParseFileSource(fileSource)
		if err != nil {
			return err
		}

		if err := addKeyFromFileToConfigMap(configmap, keyName, filePath); err != nil {
			return err
		}
	}
	return nil
}
