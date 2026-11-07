.PHONY: test coverage run

test:
	go test ./... -race -cover

coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

run:
	go run ./cmd/server
