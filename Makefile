.PHONY: build

DOCKER_NETWORK_NAME := gomiinfo-network
ENV_FILE_PATH := ./.env.json

build:
	sam build
dev: 
	sam local start-api --docker-network ${DOCKER_NETWORK_NAME} --env-vars ${ENV_FILE_PATH}
# db migration
migrate:
	migrate -path ./infra/db/migrations -database "mysql://root:password@tcp(localhost:3306)/gomi-info-db" up
migrate_down:
	migrate -path ./infra/db/migrations -database "mysql://root:password@tcp(localhost:3306)/gomi-info-db" down
# seed
import_csv: 
	sam local invoke UpdateGarbageDaysFunction --docker-network ${DOCKER_NETWORK_NAME}  --env-vars ${ENV_FILE_PATH}
lint:
	golangci-lint run ./app/...