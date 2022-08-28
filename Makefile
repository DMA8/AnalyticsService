all: docker.start test.integration docker.stop

docker.start:
	docker compose -f internal/int_tests/docker-compose.yaml up -d
	sleep 5

docker.stop:
	docker compose -f internal/int_tests/docker-compose.yaml kill

docker.restart: docker.stop docker.start

test.unit:
	go test ./... -cover

test.integration:
	go test -tags=integration ./internal/int_tests -v -count=1
