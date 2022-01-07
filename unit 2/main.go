package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	done := start(":8080")
	<-done
	log.Println("main exiting...")
}

func start(address string) <-chan struct{} {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandlerFunc)
	mux.HandleFunc("/healthz", healthzHandlerFunc)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	server := &http.Server{Addr: address, Handler: mux}
	go func() {
		<-signalChan
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("err: %v\n", err)
		} else {
			log.Printf("gracefully stopped\n")
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		log.Println("start to listen at", address, "...")

		server.ListenAndServe()
	}()

	return done
}

func defaultHandlerFunc(w http.ResponseWriter, r *http.Request) {
	for eachHeaderName, eachHeaderValue := range r.Header {
		w.Header().Set(eachHeaderName, eachHeaderValue[0])
	}

	if version := os.Getenv("VERSION"); version != "" {
		w.Header().Set("VERSION", version)
	}
	w.WriteHeader(http.StatusOK)
	log.Println("request:", r, ", ip:", r.RemoteAddr, ", response code:", http.StatusOK)
}

func healthzHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
	w.WriteHeader(http.StatusOK)
}
