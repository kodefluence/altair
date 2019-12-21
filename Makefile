test:
	@go test -cover -coverprofile=cover.out $$(go list ./... | grep -Ev "core")