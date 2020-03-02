module github.com/6RiverSystems/terraform-provider-morek8s

go 1.13

require (
	github.com/6RiverSystems/terraform-provider-helpers v0.0.2
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
	sigs.k8s.io/controller-runtime v0.5.0
)
