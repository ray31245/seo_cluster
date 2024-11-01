.PHONY publish_service_windows_release:
publish_service_windows_release:
	@echo "Publishing service for Windows release"
	GOOS=windows GOARCH=amd64 go build -o release/publish_service.exe -ldflags "\
	-X main.Url=$(URL) \
	-X main.UserName=$(USERNAME) \
	-X main.Password=$(PASSWORD) \
	" -v cmd/publish_service/main.go

.PHONY publish_service_linux_release:
publish_service_linux_release:
	@echo "Publishing service for Linux release"
	GOOS=linux GOARCH=amd64 go build -o release/publish_service -ldflags "\
	-X main.Url=$(URL) \
	-X main.UserName=$(USERNAME) \
	-X main.Password=$(PASSWORD) \
	" -v cmd/publish_service/main.go

.PHONY local_publish_service:
local_publish_service:
	@echo "Publishing service for local development"
	go build -o release/publish_service -ldflags "\
	-X main.Url=$(URL) \
	-X main.UserName=$(USERNAME) \
	-X main.Password=$(PASSWORD) \
	" -v cmd/publish_service/main.go