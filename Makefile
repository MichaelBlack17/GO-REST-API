.PHONY: build
	go build -v ./cmd/apiserver

.DEFAULT_GOAL := build