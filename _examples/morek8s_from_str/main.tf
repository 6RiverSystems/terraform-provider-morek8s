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


