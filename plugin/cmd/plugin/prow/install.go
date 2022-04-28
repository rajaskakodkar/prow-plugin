package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackageclient"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

var (
	createSecret    bool
	createConfigmap bool
	repo            = "public.ecr.aws/t0q8k6g2/repo/prow:0.1.0"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install Prow, its repo, its packages, and prerequisite secrets",
	Args:  cobra.NoArgs,
	Example: `
	tanzu prow install`,
	RunE: installProw,
}

var repoOpts = &tkgpackagedatamodel.RepositoryOptions{
	RepositoryURL:  repo,
	RepositoryName: "prow",
	Namespace:      "default",
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
	packageRepo, err := checkProwRepo(kubeConfig)
	if err == nil && packageRepo == nil {
		log.Println("Prow Repository not found.")
		log.Println("Installing Prow Repository")
		if err := installProwRepo(kubeConfig); err != nil {
			return fmt.Errorf("install prow repo: %w", err)
		}
		for {
			log.Println("Checking for Prow Repo status...")
			if packageRepo.Status.GenericStatus.FriendlyDescription == "Reconcile succeeded" {
				log.Println("Prow Repository Installed Successfully!")
				break
			} else {
				time.Sleep(10 * time.Second)
			}
		}
	} else {
		log.Println("Prow Repository exists, continuing with package installation...")
	}

	// Install packages
	installProwPackages(kubeConfig)
	return nil
}

func checkProwRepo(kubeConfig string) (*kappipkg.PackageRepository, error) {
	tkgClient, err := tkgpackageclient.NewTKGPackageClient(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("create TKG package client: %w", err)
	}
	packageRepo, err := tkgClient.GetRepository(repoOpts)
	if err != nil {
		return nil, err
	}
	return packageRepo, nil
}

func installProwRepo(kubeConfig string) error {

	tkgClient, err := tkgpackageclient.NewTKGPackageClient(kubeConfig)
	if err != nil {
		return fmt.Errorf("create TKG package client: %w", err)
	}

	progress := &tkgpackagedatamodel.PackageProgress{
		ProgressMsg: make(chan string, 10),
		Err:         make(chan error),
		Done:        make(chan struct{}),
	}

	log.Println("Adding repository")

	go tkgClient.AddRepository(repoOpts, progress, tkgpackagedatamodel.OperationTypeInstall)
	log.Println(receive(progress))
	packageRepo, err := checkProwRepo(kubeConfig)
	// todo(rajas): remove this debug line
	log.Println(packageRepo.Status.ConsecutiveReconcileSuccesses)
	return nil
}

func installProwPackages(kubeConfig string) {
	tkgClient, _ := tkgpackageclient.NewTKGPackageClient(kubeConfig)

	packages := []string{
		"crier.prow.plugin",
		"deck.prow.plugin",
		"ghproxy.prow.plugin",
		"hook.prow.plugin",
		"horologium.prow.plugin",
		"prow-cm.prow.plugin",
		"sinker.prow.plugin",
		"statusreconciler.prow.plugin",
		"tide.prow.plugin",
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

		// log.Println("Install package")
		go tkgClient.InstallPackage(packageInstallOp, progress, tkgpackagedatamodel.OperationTypeInstall)
		packageProgress := receive(progress)
		if packageProgress == nil {
			log.Println("Package Installed successfully!")
		} else {
			log.Println(packageProgress)
		}

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
