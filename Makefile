export VERSION 	?= $(shell git show -q --format=%h)
export IMAGE 		?= codefluence/altair

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

build_linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o ./build_output/linux/altair
	upx -9 -k ./build_output/linux/altair
	cp -r config/ ./build_output/linux/config/
	cp -r migration/ ./build_output/linux/migration/
	cp -r routes/ ./build_output/linux/routes/

build_darwin:
	GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -ldflags="-s" -o ./build_output/darwin/altair
	upx -9 -k ./build_output/darwin/altair
	cp -r config/ ./build_output/darwin/config/
	cp -r migration/ ./build_output/darwin/migration/
	cp -r routes/ ./build_output/darwin/routes/

build_windows:
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -ldflags="-s" -o ./build_output/windows/altair
	upx -9 -k ./build_output/windows/altair
	cp -r config/ ./build_output/windows/config/
	cp -r migration/ ./build_output/windows/migration/
	cp -r routes/ ./build_output/windows/routes/

build: build_linux build_darwin build_windows

build_docker: build_docker_latest
	sudo docker build -t $(IMAGE):$(VERSION) -f ./Dockerfile .

build_docker_latest:
	sudo docker build -t $(IMAGE):latest -f ./Dockerfile .

push_docker: push_docker_latest
	sudo docker push $(IMAGE):$(VERSION)

push_docker_latest:
	sudo docker push $(IMAGE):latest
