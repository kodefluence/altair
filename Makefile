test:
	@go test -race -cover -coverprofile=cover.out $$(go list ./... | grep -Ev "altair$$|core|mock")

mock_service:
	mockgen -source core/service.go -destination mock/mock_service.go -package mock

mock_formatter:
	mockgen -source core/formatter.go -destination mock/mock_formatter.go -package mock

mock_model:
	mockgen -source core/model.go -destination mock/mock_model.go -package mock

mock_validator:
	mockgen -source core/validator.go -destination mock/mock_validator.go -package mock

mock_plugin:
	mockgen -source core/plugin.go -destination mock/mock_plugin.go -package mock

mock_routing:
	mockgen -source core/routing.go -destination mock/mock_routing.go -package mock

mock_all: mock_service mock_formatter mock_model mock_validator mock_plugin mock_routing

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./build/linux/altair
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -o ./build/windows/altair
	GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -o ./build/darwin/altair

OPENCMD 				:=
ifeq ($(OS),Windows_NT)
	OPENCMD = start
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		OPENCMD = xdg-open
	endif
	ifeq ($(UNAME_S),Darwin)
		OPENCMD = open
	endif
endif