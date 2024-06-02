build:
	@go build -o bin/GoCRUD

run: build
	@./bin/GoCRUD

test:
	@go test -v ./...