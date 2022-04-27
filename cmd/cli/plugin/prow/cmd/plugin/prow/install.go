package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackageclient"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install Prow, its repo, its packages, and prerequisite secrets",
	Args:  cobra.NoArgs,
	Example: `
	tanzu prow install`,
	RunE: installProw,
}

// installProw will install the pro repo, its package bundles, and the
// prerequisites like secrets on a workload cluster.
func installProw(cmd *cobra.Command, _ []string) error {
	var (
		kubeConfig = getDefaultKubeconfigPath()
	)

	// Install required secrets
	createRequiredSecrets(kubeConfig)

	// Install repository
	if err := installProwRepo(kubeConfig); err != nil {
		return fmt.Errorf("install prow repo: %w", err)
	}

	// Install packages

	return nil
}

func installProwRepo(kubeconfig string) error {
	_, err := tkgpackageclient.NewTKGPackageClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("create TKG package client: %w", err)
	}

	// AddRepository validates the provided input and adds the package repository CR to the cluster
	// func (p *pkgClient) AddRepository(
	// o             *tkgpackagedatamodel.RepositoryOptions,
	// progress      *tkgpackagedatamodel.PackageProgress,
	// operationType tkgpackagedatamodel.OperationType
	// ) {
	// 	p.addRepository(o, progress, operationType)
	// }

	// tkgPkgClient.AddRepository()
	return nil
}
