SHELL=/bin/bash
UI_DIR=src/presentation/ui

dev:
	air serve

build:
	templ generate -path $(UI_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /var/infinite/ez ./ez.go

run:
	/var/infinite/ez serve