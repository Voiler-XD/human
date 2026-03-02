.PHONY: build test install clean

build:
	go build -o graphql .

test:
	go test ./... -count=1

install: build
	human plugin install --binary ./graphql

clean:
	rm -f graphql
