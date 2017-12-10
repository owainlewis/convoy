BIN := dist/convoy

all: clean build

.PHONY: build
build:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BIN} cmd/main.go

.PHONY: clean
clean:
	@rm -f dist

.PHONY: run
run:
	@go run cmd/main.go --config=config.yml --kubeconfig=/Users/owainlewis/.kube/config --v=4 --logtostderr=true

.PHONY: test
test:
	@go test -v ./...

.PHONY: image
image: clean build
	@docker build -t convoy .

