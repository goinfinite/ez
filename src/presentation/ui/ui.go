package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var frontFiles embed.FS

func UiFs() http.Handler {
	frontFileFs, err := fs.Sub(frontFiles, "dist")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(frontFileFs))
}
