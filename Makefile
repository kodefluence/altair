export VERSION 	?= $(shell git show -q --format=%h)
export IMAGE 		?= kodefluence/altair

test:
	go test -race -cover -coverprofile=cover.out $$(go list ./... | grep -Ev "altair$$|core|mock|interfaces|testhelper")

mock_metric:
	mockgen -source core/metric.go -destination mock/mock_metric.go -package mock

mock_plugin:
	mockgen -source core/plugin.go -destination mock/mock_plugin.go -package mock

mock_loader:
	mockgen -source core/cfg.go -destination mock/mock_cfg.go -package mock

mock_routing:
	mockgen -source core/routing.go -destination mock/mock_routing.go -package mock

mock_all: mock_service mock_formatter mock_model mock_validator mock_plugin mock_routing

build_linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o ./build_output/linux/altair
	upx -9 -k ./build_output/linux/altair

build_darwin:
	GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -ldflags="-s" -o ./build_output/darwin/altair
	upx -9 -k ./build_output/darwin/altair

build_windows:
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -ldflags="-s" -o ./build_output/windows/altair
	upx -9 -k ./build_output/windows/altair

build: build_linux build_darwin build_windows

build_docker: build_docker_latest
	docker build -t $(IMAGE):$(VERSION) -f ./Dockerfile .

build_docker_latest:
	docker build -t $(IMAGE):latest -f ./Dockerfile .

push_docker: push_docker_latest
	docker push $(IMAGE):$(VERSION)

tag_docker_latest:
	docker tag $(IMAGE):latest $(IMAGE):latest

push_docker_latest:
	docker push $(IMAGE):latest

docker-compose-up:
	docker-compose --env-file .env up -d

docker-compose-start:
	docker-compose --env-file .env start

docker-compose-stop:
	docker-compose stop

docker-compose-down:
	docker-compose down