build:
	docker-compose build go-rest-api

run:
	docker-compose up go-rest-api

test:
	go test -v ./...

swag:
	swag init -g cmd/app/main.go