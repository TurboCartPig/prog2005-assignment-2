.PHONY: all build test lint shell deploy

SSH_KEY = ~/.ssh/skyhigh_dennis.pem
SSH_USER = fedora
SSH_REMOTE = 10.212.142.242

all: build

build:
	go build cmd/server.go

test:
	go test ./...

lint:
	golangci-lint run

deploy: build test
	scp -i ${SSH_KEY} ./server ${SSH_USER}@${SSH_REMOTE}:/home/${SSH_USER}/server
	scp -i ${SSH_KEY} ./systemd/server.service ${SSH_USER}@${SSH_REMOTE}:/home/${SSH_USER}/.config/systemd/user/server.service
	@echo "TODO: Log onto the server and restart the service manually"
