module github.com/6RiverSystems/terraform-provider-morek8s

go 1.13

require (
	cloud.google.com/go/bigtable v1.3.0 // indirect
	github.com/6RiverSystems/terraform-provider-helpers v0.0.2
	github.com/hashicorp/terraform v0.12.21
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/terraform-providers/terraform-provider-google v1.20.0
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
	sigs.k8s.io/controller-runtime v0.5.0
)
