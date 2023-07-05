.PHONY: build

build:
	sam build
dev: 
	sam local start-api --docker-network gomiinfo-network
