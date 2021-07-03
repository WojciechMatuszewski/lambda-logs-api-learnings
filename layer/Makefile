build-LogsExtension:
	GOARCH=amd64 GOOS=linux go build -o ./extensions/logs ./extensions/logs.go
	mkdir $(ARTIFACTS_DIR)/extensions
	cp ./extensions/logs $(ARTIFACTS_DIR)/extensions/logs
