test:
	go test -race -cover -coverprofile=cover.out $$(go list ./... | grep -Ev "altair$$|core|mock|interfaces|testhelper")

mock_metric:
	mockgen -source core/metric.go -destination mock/mock_metric.go -package mock

mock_plugin:
	mockgen -source core/plugin.go -destination mock/mock_plugin.go -package mock

mock_loader:
	mockgen -source core/loader.go -destination mock/mock_loader.go -package mock

mock_routing:
	mockgen -source core/routing.go -destination mock/mock_routing.go -package mock

mock_all: mock_service mock_formatter mock_model mock_validator mock_plugin mock_routing

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./build/linux/altair
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -o ./build/windows/altair
	GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -o ./build/darwin/altair
