package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackageclient"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
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

func installProwRepo(kubeConfig string) error {
	repoOpts := &tkgpackagedatamodel.RepositoryOptions{
		RepositoryURL: "public.ecr.aws/t0q8k6g2/repo/prow@sha256:03b1bd5e1c3ec75cd66984038307db7d9dd5c2e4cea65b13ff99f2b064b3a153",
		RepositoryName: "prow",
		Namespace: "default",
	}

	progress := &tkgpackagedatamodel.PackageProgress{
		ProgressMsg: make(chan string, 10),
		Err:         make(chan error),
		Done:        make(chan struct{}),
	}


	tkgPkgClient, err := tkgpackageclient.NewTKGPackageClient(kubeConfig)
	if err != nil {
		return fmt.Errorf("create TKG package client: %w", err)
	}

	tkgPkgClient.AddRepository(repoOpts, progress, tkgpackagedatamodel.OperationTypeInstall)
	return nil
}