all: life

life: **/*.go
	go build -o life main.go

clean:
	rm -rf life
