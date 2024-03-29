# terraform-provider-morek8s

Terraform custom provider to extend standart kubernetes provider


## morek8s_from_str

Allows to manipulate with k8s resources as text in JSON or YAML format. Here is a few examples:

```
// using reference to file
resource "morek8s_from_str" "my-secret" {
  data = "${file("mysecret.yaml")}"
}

// as embedded string
resource "morek8s_from_str" "my-embedded-secret" {
  data = <<-EOF
{
    "apiVersion": "v1",
    "data": {
        "password": "MWYyZDFlMmU2N2Rm",
        "username": "YWRtaW4="
    },
    "kind": "Secret",
    "metadata": {
        "name": "my-embedded-secret",
        "namespace": "default"
    },
    "type": "Opaque"
}
  -EOF
}

```

Primary used with CRDs.

## Running Unit Tests

```
make test
```

## Running Acceptance Tests

```
make testacc
```

to run acceptance tests using your kube config:

```
KUBE_CONFIG=~/.kube/config make testacc
```

## Debugging Tests

```
KUBE_CONFIG=~/.kube/config TF_ACC=1 dlv test -- -test.run=TestAccResourceFromStr_basic -test.v
```
