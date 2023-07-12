include .env

.PHONY: build

build:
	sam build
dev: 
	sam local start-api --docker-network gomiinfo-network
migrate:
	migrate -path ./infra/db/migrations -database "mysql://root:${MYSQL_ROOT_PASSWORD}@tcp(localhost:3306)/${MYSQL_DATABASE}" up