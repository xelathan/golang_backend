build:
	@go build -o bin/golang_backend cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/golang_backend

clean:
	rm -r bin