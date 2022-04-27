package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install Prow, its repo, its packages, and prerequisite secrets",
	Args:  cobra.NoArgs,
	Example: `
	tanzu prow install`,
	RunE: installProw,
}

func installProw(cmd *cobra.Command, _ []string) error {
	fmt.Println("One day, some day soon? Install Prow, the repo, all the repo package bundles, and its prerequisites on a workload cluster.")
	return nil
}
