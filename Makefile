all: build

run:
	go run main.go

build:
	go build -o build/raylib-ai

clean:
	rm -rf build
