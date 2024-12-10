CLIENT_VERSION := $(shell cat ./version | grep "client-version" | cut -d'=' -f2)
SERVER_VERSION := $(shell cat ./version | grep "server-version" | cut -d'=' -f2)

PKGS := $(shell go list ./... 2>&1 | grep -v 'github.com/emrgen/firstime/vendor')

.PHONY: start air buf-deps proto clean-proto deps build clean lint test vet generate-client generate-docs client

start:
	go run main.go serve

init:  proto deps generate-client

info:
	@echo "client-version=$(CLIENT_VERSION)"
	@echo "server-version=$(SERVER_VERSION)"

air:
	air

# to update buf dependencies, run this command
buf-deps:
	buf mod update proto

protoc:
	@echo "Generating proto files..."
	buf generate

# build proto files
proto: clean-proto
	@echo "Generating proto files..."
	buf generate
	make generate-docs

clean-proto:
	rm -rf ./apis/v1

deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod vendor

build: protoc
	go build -o ./bin/authbase ./main/main.go

clean:
	@echo "Cleaning..."
	rm -rf vendor
	rm -rf client/generated
	rm -rf apis/v1

lint:
	@go install golang.org/x/lint/golint@latest
	for file in $(GO_FILES); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

test:
	go test  -coverprofile=profile.out -covermode=atomic $(PKGS)

vet:
	go vet $(PKGS)

generate-client: proto
	@echo "Generating client version $(CLIENT_VERSION)"
	@npx openapi-generator-cli generate \
		-i ./apis/v1/authbase.swagger.json \
		-g typescript-axios \
		-o ./clients/gen/ts/authbase-client-gen \
		--additional-properties=npmName=@emrgen/authbase-client-gen,npmVersion=${CLIENT_VERSION},useSingleRequestParameter=true,supportsES6=true,modelPropertyNaming=snake_case,paramNaming=snake_case,enumPropertyNaming=snake_case \
		--type-mappings=string=String

	# cd ./clients/firstime-gen-client/ts && yarn

generate-docs:
	@echo "Generating openapi doc version $(SERVER_VERSION)"
	@npx @redocly/cli build-docs ./apis/v1/authbase.swagger.json --output ./docs/v1/index.html

client: generate-client generate-docs
