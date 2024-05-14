build:
	docker-compose build bash_server

run:
	docker-compose up bash_server

test:
	go test cmd/main_test.go

migrate:
	migrate -path ./model -database 'postgres://postgres:password@0.0.0.0:5436/postgres?sslmode=disable' up