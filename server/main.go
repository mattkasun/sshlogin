package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPort = 8080
	timeout     = 3
	out         = time.Second * 3
)

func main() {
	port := flag.Int("p", defaultPort, "server listening port")
	flag.Usage = func() {
		fmt.Println("server which supports user registration and login via ssh")
		fmt.Printf("\nUsage:\n\n")
		fmt.Println(os.Args[0], "[flags]")
		fmt.Println("\nFlags")
		flag.PrintDefaults()
	}
	flag.Parse()
	log.Println(*port)
	router := setupRouter()
	server := http.Server{
		Addr:              fmt.Sprintf(":%d", *port),
		Handler:           router,
		ReadHeaderTimeout: out,
	}
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
		slog.Info("server shutdown")
	}()

	fmt.Println("server is listening on port", server.Addr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeout)
	defer cancel()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
		return
	}
}
