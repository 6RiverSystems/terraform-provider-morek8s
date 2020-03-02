package morek8s

import (
	"github.com/6RiverSystems/terraform-provider-helpers/kubernetes"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Provider is an entry point for resources registration
func Provider() *schema.Provider {
	providerFields := kubernetes.ProviderFields()

	p := &schema.Provider{
		Schema: providerFields,
		ResourcesMap: map[string]*schema.Resource{
			"morek8s_from_str": resourceFromStr(),
		},
	}

	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		scheme := runtime.NewScheme()
		// register needed schemas
		// apis.AddToScheme(scheme)
		// corev1.AddToScheme(scheme)

		opts := client.Options{Scheme: scheme}
		return kubernetes.NewClient(d, p.TerraformVersion, opts)
	}

	return p
}
