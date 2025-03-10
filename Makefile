# # Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# Chat support
chat-run:
	go run chat/api/service/cap/main.go | go run chat/api/tooling/logfmt/main.go

chat-test:
	curl -i -X GET http://localhost:3000/v1/status

chat-hack:
	go run chat/api/tooling/client/main.go

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor


deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor


list:
	go list -mod=mod all
