package main

import "fmt"

func createRequiredSecrets(kubeConfig string) {
	kubeClientset := getClientSet(kubeConfig)
	fmt.Println(kubeClientset)
}
