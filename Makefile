.PHONY: build

VERSION ?= 0.0.1

default: build

build:
	go build -o terraform-provider-morek8s_v$(VERSION)

init: build
	terraform init

