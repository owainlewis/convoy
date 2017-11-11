.PHONY: build
build:
	go build cmd/main.go

.PHONY: run
run:
	@go run cmd/main.go --config=/Users/owainlewis/.kube/config
