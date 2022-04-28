package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackageclient"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

var (
	createSecret    bool
	createConfigmap bool
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install Prow, its repo, its packages, and prerequisite secrets",
	Args:  cobra.NoArgs,
	Example: `
	tanzu prow install`,
	RunE: installProw,
}

func init() {
	InstallCmd.Flags().BoolVarP(&createSecret, "create-secrets", "", false, "Should create secrets")
	InstallCmd.Flags().BoolVarP(&createConfigmap, "create-configmaps", "", false, "Should create configmaps")
}

// installProw will install the pro repo, its package bundles, and the
// prerequisites like secrets on a workload cluster.
func installProw(cmd *cobra.Command, _ []string) error {
	var (
		kubeConfig = getDefaultKubeconfigPath()
	)

	if createSecret {
		// Install required secrets
		err := createRequiredSecrets(kubeConfig)
		if err != nil {
			panic(err)
		}
	}

	if createConfigmap {
		// Install required secrets
		err := createRequiredConfigmap(kubeConfig)
		if err != nil {
			panic(err)
		}
	}

	// Install repository
	/*if err := installProwRepo(kubeConfig); err != nil {
		return fmt.Errorf("install prow repo: %w", err)
	}*/
	// Install packages
	installProwPackages(kubeConfig)
	return nil
}

func installProwRepo(kubeConfig string) error {
	repo := "public.ecr.aws/t0q8k6g2/repo/prow@sha256:03b1bd5e1c3ec75cd66984038307db7d9dd5c2e4cea65b13ff99f2b064b3a153"

	tkgClient, err := tkgpackageclient.NewTKGPackageClient(kubeConfig)
	if err != nil {
		return fmt.Errorf("create TKG package client: %w", err)
	}

	repoOpts := &tkgpackagedatamodel.RepositoryOptions{
		RepositoryURL:  repo,
		RepositoryName: "prow",
		Namespace:      "default",
	}

	progress := &tkgpackagedatamodel.PackageProgress{
		ProgressMsg: make(chan string, 10),
		Err:         make(chan error),
		Done:        make(chan struct{}),
	}

	log.Println("Adding repository")
	go tkgClient.AddRepository(repoOpts, progress, tkgpackagedatamodel.OperationTypeInstall)
	log.Println(receive(progress))

	return nil
}

func installProwPackages(kubeConfig string) {
	tkgClient, _ := tkgpackageclient.NewTKGPackageClient(kubeConfig)

	packages := []string{
		"crier.prow.plugin",
	}

	for _, pkg := range packages {
		log.Printf("Installing package: %v\n", pkg)

		var packageInstallOp = tkgpackagedatamodel.NewPackageOptions()
		packageInstallOp.PkgInstallName = pkg
		packageInstallOp.PackageName = pkg
		packageInstallOp.Namespace = "default"
		packageInstallOp.Version = "0.1.0"

		progress := &tkgpackagedatamodel.PackageProgress{
			ProgressMsg: make(chan string, 10),
			Err:         make(chan error),
			Done:        make(chan struct{}),
		}

		log.Println("Install package")
		go tkgClient.InstallPackage(packageInstallOp, progress, tkgpackagedatamodel.OperationTypeInstall)
		log.Println(receive(progress))
	}
}
func receive(progress *tkgpackagedatamodel.PackageProgress) error {
	for {
		select {
		case err := <-progress.Err:
			log.Println("ERROR")
			return err
		case msg := <-progress.ProgressMsg:
			log.Println(msg)
			continue
		case <-progress.Done:
			return nil
		}
	}
}
