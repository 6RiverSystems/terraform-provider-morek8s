.PHONY: build

VERSION ?= 0.0.1
TEST?="./morek8s"

default: build

test:
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

build:
	go build -o terraform-provider-morek8s_v$(VERSION)

init: build
	terraform init

