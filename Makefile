.PHONY: build

build:
	sam build
dev: 
	sam local start-api --docker-network gomiinfo-network --env-vars ./.env.json
migrate:
	migrate -path ./infra/db/migrations -database "mysql://root:password@tcp(localhost:3306)/gomi-info-db" up
seed: 
	sam local invoke UpdateGarbageDaysFunction
lint:
	golangci-lint run ./app/...