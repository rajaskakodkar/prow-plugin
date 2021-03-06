package main

import (
	"os"

	"github.com/aunum/log"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
)

var descriptor = cliv1alpha1.PluginDescriptor{
	Name:        "prow",
	Description: "Tanzupalooza 2022 Prow Plugin",
	Version:     "v0.0.1",
	Group:       cliv1alpha1.ManageCmdGroup, // set group
}

func main() {
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		log.Fatal(err)
	}
	p.AddCommands(
		InstallCmd,
	)
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
