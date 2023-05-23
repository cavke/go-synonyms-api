APP?=synonym-api

# -----------------------------------------------------------------
#				Main targets
# -----------------------------------------------------------------

.PHONY: build
## build: builds the application
build:
	go build -o ${APP} cmd/main.go

.PHONY: run
## run: runs go run main.go
run: build
	go run cmd/main.go

.PHONY: run-port
## run: runs go run main.go
run-port: build
	go run cmd/main.go -port=${port}

.PHONY: clean
## clean: cleans the binary
clean:
	go clean

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...

.PHONY: lint
## lint: runs linter
lint:
	golangci-lint run

# -----------------------------------------------------------------
#				Docker targets
# -----------------------------------------------------------------

.PHONY: docker-build
## docker-build: builds the app docker image to registry
docker-build:
	docker build -t ${APP} -f deploy/dockerfile .

.PHONY: docker-run
## docker-run: starts the program instances
docker-run:
	docker run -p 8080:8080 ${APP}

# -----------------------------------------------------------------
#			 Help
# -----------------------------------------------------------------
.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |	sed -e 's/^/ /'

