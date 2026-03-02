.PHONY: build test install clean

build:
	go build -o k8s .

test:
	go test ./... -count=1

install: build
	human plugin install --binary ./k8s

clean:
	rm -f k8s
