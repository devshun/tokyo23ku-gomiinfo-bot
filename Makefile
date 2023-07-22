.PHONY: build

DOCKER_NETWORK_NAME := gomi-info-network
ENV_FILE_PATH := ./.env.json

build:
	rm -rf ./.aws-sam ; sam build
dev: 
	sam local start-api --docker-network ${DOCKER_NETWORK_NAME} --env-vars ${ENV_FILE_PATH}
# db migration
migrate:
	migrate -path ./db/migrations -database "mysql://root:password@tcp(localhost:3306)/gomi-info-db" up
migrate_down:
	migrate -path ./db/migrations -database "mysql://root:password@tcp(localhost:3306)/gomi-info-db" down
# import csv
import_csv: 
	sam local invoke ImportCsvToDbFunction --docker-network ${DOCKER_NETWORK_NAME}  --env-vars ${ENV_FILE_PATH}
# lint
lint:
	cd app && golangci-lint run ./...
# fmt
fmt:
	cd app && go fmt ./...
# test
test:
	cd app && go test ./...