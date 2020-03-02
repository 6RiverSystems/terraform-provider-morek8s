package main

import (
	"github.com/6RiverSystems/terraform-provider-morek8s/morek8s"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return morek8s.Provider()
		},
	})
}
