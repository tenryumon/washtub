docker-start: deps
	@echo "Building and Starting Container..."
	@docker-compose -f docker/docker-compose.yml up -d

docker-stop: deps
	@echo "Stopping Container..."
	@docker-compose -f docker/docker-compose.yml stop

docker-remove: docker-stop
	@echo "Removing Container..."
	@docker-compose -f docker/docker-compose.yml down

run:
	@echo "Running Service..."
	@go run main.go

test:
	@echo "Unit-Testing..."
	@go test ./... -race -cover


