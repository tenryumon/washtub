DOCKER := $(shell which docker)
DOCKER_COMPOSE := $(shell which docker-compose)
GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
BINARY_NAME=serviceapa
BINARY_UNIX=$(BINARY_NAME)_unix

APP_NAME := backend

DOCKER_IMAGE_NAME := $(APP_NAME)
ifdef AWS_ECR_REGISTRY
	DOCKER_IMAGE_NAME := $(addprefix "${AWS_ECR_REGISTRY}/", $(DOCKER_IMAGE_NAME))
endif

# GIT_COMMIT is the commit hash for a specific commit in a branch. Please NOTE that the length of the git-commit
# is truncated to 14 characters.
GIT_COMMIT := $(shell git rev-parse HEAD | cut -c 1-14)
ifdef RELEASE_TAG
	GIT_COMMIT := ${RELEASE_TAG}
endif

deps:
ifeq (,$(DOCKER))
	@echo "docker is not available. Please install docker in https://docs.docker.com/desktop/"
	@exit 1
endif
ifeq (,$(DOCKER_COMPOSE))
	@echo "docker-compose is not available. Please install docker-compose in https://docs.docker.com/desktop/"
	@exit 1
endif

docker-start: deps
	@echo "Building and Starting Container..."
	@docker-compose -f docker/docker-compose.yml up -d

docker-stop: deps
	@echo "Stopping Container..."
	@docker-compose -f docker/docker-compose.yml stop

docker-remove: deps docker-stop
	@echo "Removing Container..."
	@docker-compose -f docker/docker-compose.yml down

docker-list: deps
	@sudo docker ps -a

run:
	@echo "Running Service..."
	@go run cmd/api/main.go

test:
	@echo "Unit-Testing..."
	@go test ./... -race -cover -p=1

.PHONY: build
build:
	CGO_ENABLED=0 go build -o backend  -v cmd/api/*.go

# parse-config parse the configuration based on the $ENV so we can inject the configuration
# later on for deployment.
.PHONY: parse-config
parse-config:
	go run cmd/tools/config_parse.go \
		-config-path=devops/configuration/backend-${ENV}.ini \
		-target-path=backend-${ENV}.ini

run-db:
	@echo "Do Changes in Database using sqlfile with prefix ${prefix}"
	@go run cmd/tools/db_populator.go -all -sqlprefix=${prefix}

.PHONY: seed-test
seed-test:
	@echo "Seeding Test Data..."
	@go run cmd/tools/db_populator.go -all -sqlprefix=${prefix} -sqlfolder=sqlseed/local/

build-linux:
	@CGO_ENABLED=0 GOOS=linux go build -o unix -v cmd/api/main.go

.PHONY: docker-build
docker-build: build parse-config
	@cp devops/dockerfiles/go-x86_64 Dockerfile
	@sed -i -- 's/envname/${ENV}/g; s/awsid/${AWS_ACCESS_KEY_ID}/g; s/awssecret/${AWS_SECRET_ACCESS_KEY}/g' Dockerfile
	@docker buildx build --platform linux/amd64 \
		-f Dockerfile \
		-t $(DOCKER_IMAGE_NAME):${GIT_COMMIT} \
		-t $(DOCKER_IMAGE_NAME):latest \
		--load .
	@rm backend
	@rm Dockerfile
	@rm backend-${ENV}.ini
	
.PHONY: docker-push
docker-push:
ifdef AWS_ECR_REGISTRY
	@aws ecr describe-repositories --repository-names $(APP_NAME) || aws ecr create-repository --repository-name $(APP_NAME)
endif
	@docker image push --all-tags $(DOCKER_IMAGE_NAME)
