package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/two-hundred/celerity-github-registry/internal/registry"
)

func main() {
	router := mux.NewRouter()
	port, accessLogWriter, err := registry.Setup(
		router,
		registry.GetDependencies,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer accessLogWriter.Close()

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		// Same as the ALB default idle timeout.
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           http.TimeoutHandler(router, 60*time.Second, "Timeout!\n"),
	}
	log.Printf("Starting server on port %d", port)
	log.Fatal(srv.ListenAndServe())
}
