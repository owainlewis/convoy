all: build

.PHONY: build
build:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build cmd/main.go

.PHONY: run
run:
	@go run cmd/main.go --config=/Users/owainlewis/.kube/config --v=4 --logtostderr=true
