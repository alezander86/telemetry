.PHONY: build test

build:
	GOOS=linux GOARCH=amd64 sam build --no-cached

build-go:
	cd telemetry-collector;GOOS=linux GOARCH=amd64 go build -o ../bin/telemetry-collector

test:
	cd telemetry-collector;go test ./... -coverprofile=coverage.out `go list ./...`

deploy:
	sam deploy --no-confirm-changeset --no-fail-on-empty-changeset

start-local-api: build
	sam local start-api
