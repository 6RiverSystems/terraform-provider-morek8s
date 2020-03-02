#!/bin/bash

GOOS=darwin GOARCH=amd64 go build -v -o build/terraform-provider-morek8s_darwin-amd64
GOOS=linux GOARCH=amd64 go build -v -o build/terraform-provider-morek8s_linux-amd64
GOOS=windows GOARCH=amd64 go build -v -o build/terraform-provider-morek8s_windows-amd64

gzip build/*
