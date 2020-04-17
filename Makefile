test:
	@go test -cover -coverprofile=cover.out $$(go list ./... | grep -Ev "altair$$|core|mock")

mock_service:
	mockgen -source core/service.go -destination mock/mock_service.go -package mock

mock_formatter:
	mockgen -source core/formatter.go -destination mock/mock_formatter.go -package mock

mock_model:
	mockgen -source core/model.go -destination mock/mock_model.go -package mock

mock_validator:
	mockgen -source core/validator.go -destination mock/mock_validator.go -package mock

mock_routing:
	mockgen -source core/routing.go -destination mock/mock_routing.go -package mock

mock_all: mock_service mock_formatter mock_model mock_validator

generate_blueprint:
	snowboard apib -o blueprint/_output/API.apib blueprint/API.md
	snowboard html -o blueprint/_output/index.html blueprint/_output/API.apib

open_blueprint:
	$(OPENCMD) blueprint/_output/index.html

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