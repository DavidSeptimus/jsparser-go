.PHONY: run docker-build docker-run docker


all: docker


# runs jsparser-go in the local golang environment
run:
	go run main.go

# builds the jsparser-go docker image
docker-build:
	docker compose build

# executes jsparser-go in a docker container
docker-run:
	docker compose up

# both build and run jsparser-go
docker: docker-build docker-run

# destroys the compose app (container + network
down:
	docker compose down
