all: tidy build run

tidy:
	go mod tidy

build:
	go build -o bin/weather cmd/main.go

run:
	./bin/weather

clean:
	rm -rf bin/weather
