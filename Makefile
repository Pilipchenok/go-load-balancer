.PHONY: build run test race vet clean docker-up docker-down

build:
	go build -o bin/balancer ./cmd/balancer

run:
	go run ./cmd/balancer

test:
	go test ./... -v

race:
	go test -race ./... -v

vet:
	go vet ./...

clean:
	rm -rf bin/

docker-up:
	docker compose up --build

docker-down:
	docker compose down