.PHONY: build run test race lint clean release

build:
	go build -trimpath -ldflags="-s -w" -o bin/soundcloud-radio ./cmd/radio

run: build
	./bin/soundcloud-radio

test:
	go test ./...

race:
	go test -race ./...

lint:
	go vet ./...

clean:
	rm -rf bin/

release:
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o bin/soundcloud-radio-linux-amd64 ./cmd/radio
	GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o bin/soundcloud-radio-linux-arm64 ./cmd/radio
