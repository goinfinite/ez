SHELL=/bin/bash
UI_DIR=src/presentation/ui

dev:
	air serve

refresh:
	# TODO: Refresh frontend using a websocket

build:
	templ generate -path $(UI_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /var/speedia/control ./control.go

run:
	/var/speedia/control serve