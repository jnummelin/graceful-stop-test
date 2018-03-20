package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var healthy bool

func hello(w http.ResponseWriter, r *http.Request) {
	time.Sleep(500 * time.Millisecond)
	fmt.Fprintf(w, "Hello World!")
}

func ping(w http.ResponseWriter, r *http.Request) {
	if !healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	fmt.Fprintf(w, "pong")
}

func main() {

	healthy = true

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/ping", ping)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("Starting to listen on :8080")
		log.Fatal(s.ListenAndServe())
	}()

	// Put in a signal handler for SIGTERM
	log.Println("Setting signal handler")
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	<-signals
	log.Println("Got shutdown signal, waiting 15s to finish ongoing reqs")
	// This will make the healthcheck fail --> LBs will not give any more traffic
	healthy = false

	// Wait for the ongoing requests to finish
	// Note: s.Shutdown actually does gracefull shutdown, here we just make things bit more explicit for
	// demonstration purposes
	time.Sleep(10 * time.Second)
	// Stop the server
	s.Shutdown(nil)
}
