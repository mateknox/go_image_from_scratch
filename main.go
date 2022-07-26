package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
)

//include whole static dir
//go:embed static
var content embed.FS

func handler() http.Handler {
	fsys := fs.FS(content)
	html, _ := fs.Sub(fsys, "static")
	return http.FileServer(http.FS(html))
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", handler())
	err := http.ListenAndServe(":5555", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("type: %T; value: %q\n", err, err)
		os.Exit(1)
	}
}
