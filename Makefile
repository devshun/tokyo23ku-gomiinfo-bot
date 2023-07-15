.PHONY: build

DOCKER_NETWORK := gomiinfo-network
ENV_FILE_PATH := ./.env.json

build:
	sam build
dev: 
	sam local start-api --docker-network ${DOCKER_NETWORK} --env-vars ${ENV_FILE_PATH}
migrate:
	migrate -path ./infra/db/migrations -database "mysql://root:password@tcp(localhost:3306)/gomi-info-db" up
seed: 
	sam local invoke UpdateGarbageDaysFunction --docker-network ${DOCKER_NETWORK}  --env-vars ${ENV_FILE_PATH}
lint:
	golangci-lint run ./app/...