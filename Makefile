.PHONY: cover viewcover lint

cover:
	go test -v -race -coverpkg=./... -coverprofile=coverage.out ./...

viewcover:
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./...
