package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Include the entire static assets directory inside the compiled binary footprint
//go:embed static
var embeddedContent embed.FS

func getStaticFileSystem() http.FileSystem {
	// Isolate the "static" subdirectory from the embedded engine root
	subFS, err := fs.Sub(embeddedContent, "static")
	if err != nil {
		log.Fatalf("Critical error initializing embedded asset path: %v", err)
	}
	return http.FS(subFS)
}

func main() {
	// Initialize standard routing multiplexer
	mux := http.NewServeMux()

	// Route all traffic root paths directly into our file system
	mux.Handle("/", http.FileServer(getStaticFileSystem()))

	server := &http.Server{
		Addr:         ":5555",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Channel to capture incoming termination signals safely
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Spin up the web listener loop asynchronously inside a goroutine
	go func() {
		log.Printf("Microservice engine online. Listening on interface 0.0.0.0:5555")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Critical listener failure: %v", err)
		}
	}()

	// Execution blocks here until a termination signal is intercepted
	sig := <-shutdownChan
	log.Printf("Termination signal received (%v). Starting graceful shutdown...", sig)

	// Allow open web connections up to 5 seconds to finish their requests safely
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown abruptly: %v", err)
	} else {
		log.Println("Server exited cleanly.")
	}
}
