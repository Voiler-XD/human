.PHONY: build test install clean

build:
	go build -o htmx .

test:
	go test ./... -count=1

install: build
	human plugin install --binary ./htmx

clean:
	rm -f htmx
