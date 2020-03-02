.PHONY: build

default: build

build:
	go build -o terraform-provider-morek8s

init: build
	terraform init

