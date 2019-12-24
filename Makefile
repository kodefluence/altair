test:
	@go test -cover -coverprofile=cover.out $$(go list ./... | grep -Ev "altair$$|core|mock")

mock_service:
	mockgen -source core/service.go -destination mock/mock_service.go -package mock

mock_formatter:
	mockgen -source core/formatter.go -destination mock/mock_formatter.go -package mock

mock_model:
	mockgen -source core/model.go -destination mock/mock_model.go -package mock

mock_all: mock_service mock_formatter mock_model