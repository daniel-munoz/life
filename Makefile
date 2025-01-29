all: life

life: **/*.go
	go build -o life main.go

test: **/*.go
	go test -v ./...

coverage: **/*.go
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf life coverage.out
