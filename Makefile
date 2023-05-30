SHELL := /bin/bash

build:
	@docker build -t http-benchmarking .

run:
	@make build && docker run -it http-benchmarking

prometheus:
	@rm -rf ./k8s/deploy && ./k8s.sh

apply:
	@kubectl apply -f ./k8s/deploy

.PHONY: k8s
