db-up:
	docker run -d --name postgres_container -p 5432:5432 \
		-e POSTGRES_PASSWORD=postgres \
		-e DB_HOST=localhost -e DB_PORT=5432 -e DB_NAME=postgres \
		-e DB_USER=postgres -e DB_PASS=postgres \
		postgres:16.2

test:
	go test ./...
	
test-coverage:
	go test -coverprofile=coverage.out ./...

coverage: test
	go tool cover -html=coverage.out -o coverage.html

clean-db:
	docker stop $$(docker ps -q --filter ancestor=postgres:16.2) || true
	docker rm $$(docker ps -aq --filter ancestor=postgres:16.2) || true

db-login:
	docker exec -it postgres_container psql -U postgres -d postgres

run:
	go run cmd/shell/main.go
